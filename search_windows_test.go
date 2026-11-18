package main

import (
	"strings"
	"syscall"
	"testing"
)

func TestSearchReportsWindowsSharingViolation(t *testing.T) {
	binary := buildCLI(t)
	root := t.TempDir()
	writeFixture(t, root, "a.txt", "needle in a readable file\n")
	locked := writeFixture(t, root, "b.txt", "needle in a locked file\n")
	name, err := syscall.UTF16PtrFromString(locked)
	if err != nil {
		t.Fatal(err)
	}
	handle, err := syscall.CreateFile(name, syscall.GENERIC_READ, 0, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		t.Fatalf("lock fixture file: %v", err)
	}
	t.Cleanup(func() {
		if err := syscall.CloseHandle(handle); err != nil {
			t.Errorf("close fixture handle: %v", err)
		}
	})

	result := runCLI(t, binary, root, "search", "needle")
	assertCLI(t, result, 2, "a.txt:1:needle in a readable file\n", true)
	if !strings.Contains(result.stderr, "b.txt") {
		t.Errorf("diagnostic %q does not identify b.txt", result.stderr)
	}
}
