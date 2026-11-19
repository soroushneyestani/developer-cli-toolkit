package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
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
		valid, err := collectMatches(reader, file, literal, &matches)
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

func collectMatches(reader *bufio.Reader, source io.ReaderAt, literal []byte, matches *stagedMatches) (bool, error) {
	var offset int64
	var stagingErr error
	for lineNumber := 1; ; lineNumber++ {
		line, err := reader.ReadSlice('\n')
		if err == bufio.ErrBufferFull {
			long, err := scanLongLine(reader, line, literal)
			if err != nil && err != io.EOF {
				return false, err
			}
			if !long.valid {
				return false, drainRejectedSource(reader, err)
			}
			if stagingErr == nil && long.matched {
				if err := matches.stageLongLine(lineNumber, source, offset, long); err != nil {
					stagingErr = fmt.Errorf("stage search matches: %w", err)
				}
			}
			offset += long.consumed
			if err == io.EOF {
				return true, stagingErr
			}
			continue
		}
		if err != nil && err != io.EOF {
			return false, err
		}
		if len(line) == 0 && err == io.EOF {
			return true, stagingErr
		}
		if bytes.IndexByte(line, 0) >= 0 || !utf8.Valid(line) {
			return false, drainRejectedSource(reader, err)
		}
		offset += int64(len(line))
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

func drainRejectedSource(reader *bufio.Reader, lastRead error) error {
	// A read failure makes the command fail even for a file we would skip.
	if lastRead == io.EOF {
		return nil
	}
	_, err := io.Copy(io.Discard, reader)
	return err
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
	s.encodeHeader(lineNumber, int64(len(line)))
	if memory, ok := s.reserveMemory(int64(len(line))); ok {
		copy(memory, line)
		return nil
	}
	if err := s.startSpill(); err != nil {
		return err
	}
	if _, err := s.writer.Write(s.header[:]); err != nil {
		return err
	}
	_, err := s.writer.Write(line)
	return err
}

func (s *stagedMatches) stageLongLine(lineNumber int, source io.ReaderAt, offset int64, line longLineResult) error {
	s.encodeHeader(lineNumber, line.size)
	section := io.NewSectionReader(source, offset, line.size)
	if memory, ok := s.reserveMemory(line.size); ok {
		_, err := io.ReadFull(section, memory)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}
		if err != nil || sha256.Sum256(memory) != line.digest {
			return fmt.Errorf("source changed while staging long line %d", lineNumber)
		}
		return nil
	}
	if err := s.startSpill(); err != nil {
		return err
	}
	if _, err := s.writer.Write(s.header[:]); err != nil {
		return err
	}
	digest := sha256.New()
	copied, err := io.Copy(io.MultiWriter(s.writer, digest), section)
	if err != nil {
		return err
	}
	var actual [sha256.Size]byte
	digest.Sum(actual[:0])
	if copied != line.size || actual != line.digest {
		return fmt.Errorf("source changed while staging long line %d", lineNumber)
	}
	return nil
}

func (s *stagedMatches) reserveMemory(width int64) ([]byte, bool) {
	if s.spill != nil || width > int64(matchMemoryLimit-len(s.memory)-matchHeaderSize) {
		return nil, false
	}
	start := len(s.memory)
	needed := start + matchHeaderSize + int(width)
	if needed > cap(s.memory) {
		larger := make([]byte, start, min(matchMemoryLimit, max(needed, cap(s.memory)*2)))
		copy(larger, s.memory)
		s.memory = larger
	}
	s.memory = s.memory[:needed]
	copy(s.memory[start:], s.header[:])
	return s.memory[start+matchHeaderSize:], true
}

func (s *stagedMatches) encodeHeader(lineNumber int, width int64) {
	binary.LittleEndian.PutUint64(s.header[:8], uint64(lineNumber))
	binary.LittleEndian.PutUint64(s.header[8:], uint64(width))
}

func (s *stagedMatches) startSpill() error {
	if s.spill != nil {
		return nil
	}
	file, err := os.CreateTemp("", "dct-matches-*")
	if err != nil {
		return err
	}
	s.spill = file
	s.writer = bufio.NewWriterSize(file, searchBufferSize)
	_, err = s.writer.Write(s.memory)
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
		var text string
		if width > searchBufferSize {
			var builder strings.Builder
			builder.Grow(width)
			if _, err := io.CopyN(&builder, input, int64(width)); err != nil {
				return found, fmt.Errorf("read temporary search match text: %w", err)
			}
			text = builder.String()
		} else {
			line = slices.Grow(line[:0], width)[:width]
			if _, err := io.ReadFull(input, line); err != nil {
				return found, fmt.Errorf("read temporary search match text: %w", err)
			}
			text = string(line)
		}
		if err := emit(searchMatch{Path: path, Line: lineNumber, Text: text}); err != nil {
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
