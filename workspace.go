package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"
)

func openDirectory(path string) (*os.Root, error) {
	path = filepath.Clean(path)
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("root %q must be a directory and cannot be a symbolic link", path)
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		return nil, err
	}
	opened, err := root.Stat(".")
	if err != nil {
		root.Close()
		return nil, err
	}
	if !os.SameFile(info, opened) {
		root.Close()
		return nil, fmt.Errorf("root %q changed while being opened", path)
	}
	return root, nil
}

func listFiles(root *os.Root, extensions []string) ([]string, error) {
	var files []string
	err := fs.WalkDir(root.FS(), ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		if len(extensions) > 0 && !slices.ContainsFunc(extensions, func(extension string) bool {
			return strings.HasSuffix(entry.Name(), extension)
		}) {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Traversal order can differ from the order of complete slash-separated paths.
	slices.Sort(files)
	return files, nil
}

type searchMatch struct {
	Path string `json:"path"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

func searchFiles(root *os.Root, files []string, literal []byte, emit func(searchMatch) error) (bool, error) {
	found := false
	for _, file := range files {
		contents, err := root.ReadFile(file)
		if err != nil {
			return found, err
		}
		// Validate the whole file before emitting any matches from it.
		if bytes.IndexByte(contents, 0) >= 0 || !utf8.Valid(contents) {
			continue
		}
		for lineNumber := 1; len(contents) > 0; lineNumber++ {
			line, rest, terminated := bytes.Cut(contents, []byte{'\n'})
			contents = rest
			if terminated {
				line = bytes.TrimSuffix(line, []byte{'\r'})
			}
			if bytes.Contains(line, literal) {
				if err := emit(searchMatch{Path: file, Line: lineNumber, Text: string(line)}); err != nil {
					return found, fmt.Errorf("write search result: %w", err)
				}
				found = true
			}
		}
	}
	return found, nil
}

type fileSize struct {
	Path  string `json:"path"`
	Bytes int64  `json:"bytes"`
}

func inspectFiles(root *os.Root, files []string, emit func(fileSize) error) error {
	for _, file := range files {
		info, err := root.Lstat(file)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			continue
		}
		if err := emit(fileSize{Path: file, Bytes: info.Size()}); err != nil {
			return fmt.Errorf("write inspection result: %w", err)
		}
	}
	return nil
}
