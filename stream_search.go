package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"unicode/utf8"
)

const (
	searchBufferSize = 64 * 1024
	matchMemoryLimit = 256 * 1024
	matchHeaderSize  = 16
)

type searchMatch struct {
	Path string `json:"path"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

func searchFiles(root *os.Root, files []string, literal []byte, emit func(searchMatch) error) (bool, error) {
	reader := bufio.NewReaderSize(nil, searchBufferSize)
	var matches stagedMatches
	found := false
	for _, path := range files {
		file, err := root.Open(path)
		if err != nil {
			return found, err
		}
		reader.Reset(file)
		valid, err := collectMatches(reader, literal, &matches)
		file.Close()

		if valid && err == nil {
			var emitted bool
			emitted, err = matches.replay(path, reader, emit)
			found = found || emitted
		}
		cleanupErr := matches.reset()
		if err != nil {
			return found, err
		}
		if cleanupErr != nil {
			return found, fmt.Errorf("remove temporary search matches: %w", cleanupErr)
		}
	}
	return found, nil
}

func collectMatches(reader *bufio.Reader, literal []byte, matches *stagedMatches) (bool, error) {
	var longLine []byte
	var stagingErr error
	for lineNumber := 1; ; lineNumber++ {
		line, err := reader.ReadSlice('\n')
		if err == bufio.ErrBufferFull {
			longLine = append(longLine[:0], line...)
			for err == bufio.ErrBufferFull {
				line, err = reader.ReadSlice('\n')
				longLine = append(longLine, line...)
			}
			line = longLine
		}
		if err != nil && err != io.EOF {
			return false, err
		}
		if len(line) == 0 && err == io.EOF {
			return true, stagingErr
		}
		if bytes.IndexByte(line, 0) >= 0 || !utf8.Valid(line) {
			// A read failure makes the command fail even for a file we would skip.
			if err == io.EOF {
				return false, nil
			}
			_, err = io.Copy(io.Discard, reader)
			return false, err
		}
		if err == nil {
			line = line[:len(line)-1]
			line = bytes.TrimSuffix(line, []byte{'\r'})
		}
		if stagingErr == nil && bytes.Contains(line, literal) {
			// The staging error matters only if the source is valid and readable.
			if err := matches.add(lineNumber, line); err != nil {
				stagingErr = fmt.Errorf("stage search matches: %w", err)
			}
		}
		if err == io.EOF {
			return true, stagingErr
		}
	}
}

// stagedMatches preserves full-file validation without retaining all matching
// text in memory. Each record stores its line number, byte length, and text.
type stagedMatches struct {
	memory []byte
	spill  *os.File
	writer *bufio.Writer
	header [matchHeaderSize]byte
}

func (s *stagedMatches) add(lineNumber int, line []byte) error {
	binary.LittleEndian.PutUint64(s.header[:8], uint64(lineNumber))
	binary.LittleEndian.PutUint64(s.header[8:], uint64(len(line)))
	if s.spill == nil && len(line) <= matchMemoryLimit-len(s.memory)-matchHeaderSize {
		needed := len(s.memory) + matchHeaderSize + len(line)
		if needed > cap(s.memory) {
			larger := make([]byte, len(s.memory), min(matchMemoryLimit, max(needed, cap(s.memory)*2)))
			copy(larger, s.memory)
			s.memory = larger
		}
		s.memory = append(s.memory, s.header[:]...)
		s.memory = append(s.memory, line...)
		return nil
	}
	if s.spill == nil {
		file, err := os.CreateTemp("", "dct-matches-*")
		if err != nil {
			return err
		}
		s.spill = file
		s.writer = bufio.NewWriterSize(file, searchBufferSize)
		if _, err := s.writer.Write(s.memory); err != nil {
			return err
		}
	}
	if _, err := s.writer.Write(s.header[:]); err != nil {
		return err
	}
	_, err := s.writer.Write(line)
	return err
}

func (s *stagedMatches) replay(path string, buffer *bufio.Reader, emit func(searchMatch) error) (bool, error) {
	if s.spill == nil && len(s.memory) == 0 {
		return false, nil
	}
	var input io.Reader = bytes.NewReader(s.memory)
	if s.spill != nil {
		if err := s.writer.Flush(); err != nil {
			return false, fmt.Errorf("stage search matches: %w", err)
		}
		if _, err := s.spill.Seek(0, io.SeekStart); err != nil {
			return false, fmt.Errorf("rewind temporary search matches: %w", err)
		}
		buffer.Reset(s.spill)
		input = buffer
	}

	found := false
	var line []byte
	for {
		if _, err := io.ReadFull(input, s.header[:]); err != nil {
			if err == io.EOF {
				return found, nil
			}
			return found, fmt.Errorf("read temporary search match header: %w", err)
		}
		lineNumber := int(binary.LittleEndian.Uint64(s.header[:8]))
		width := int(binary.LittleEndian.Uint64(s.header[8:]))
		line = slices.Grow(line[:0], width)[:width]
		if _, err := io.ReadFull(input, line); err != nil {
			return found, fmt.Errorf("read temporary search match text: %w", err)
		}
		if err := emit(searchMatch{Path: path, Line: lineNumber, Text: string(line)}); err != nil {
			return found, fmt.Errorf("write search result: %w", err)
		}
		found = true
	}
}

func (s *stagedMatches) reset() error {
	s.memory = s.memory[:0]
	if s.spill == nil {
		return nil
	}
	err := errors.Join(s.spill.Close(), os.Remove(s.spill.Name()))
	s.spill = nil
	s.writer = nil
	return err
}
