package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
