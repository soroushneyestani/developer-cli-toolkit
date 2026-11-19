package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"
)

func TestCollectMatchesReturnsLateSourceReadErrors(t *testing.T) {
	failure := errors.New("test source read failure")
	for _, tail := range []string{"valid text\n", "\x00\n", "\xff\n"} {
		source := io.MultiReader(
			strings.NewReader("needle before error\n"+tail+"data beyond first buffer\n"),
			iotest.ErrReader(failure),
		)
		var staged stagedMatches
		valid, err := collectMatches(bufio.NewReaderSize(source, 16), []byte("needle"), &staged)
		if valid || !errors.Is(err, failure) {
			t.Errorf("source with tail %q: valid %t, error %v, want source read failure", tail, valid, err)
		}
	}
}

func TestCollectMatchesDefersStagingErrorUntilValidatedEOF(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "closed-spill-*")
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	for _, valid := range []bool{true, false} {
		text := "needle\n"
		if !valid {
			text += "\x00\n"
		}
		staged := stagedMatches{spill: file, writer: bufio.NewWriterSize(file, 1)}
		gotValid, err := collectMatches(bufio.NewReader(strings.NewReader(text)), []byte("needle"), &staged)
		if valid != gotValid || (valid && !errors.Is(err, os.ErrClosed)) || (!valid && err != nil) {
			t.Errorf("valid source %t: valid %t, error %v", valid, gotValid, err)
		}
	}
}

func TestSearchFilesRemovesSpill(t *testing.T) {
	writeFailure := errors.New("test search output failure")
	longLine := "needle " + strings.Repeat("x", matchMemoryLimit)
	cases := []struct {
		name     string
		tail     string
		writeErr error
		visits   int
		found    bool
	}{
		{name: "valid file", visits: 2, found: true},
		{name: "invalid file", tail: "\x00\n"},
		{name: "output error", writeErr: writeFailure, visits: 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rootPath := t.TempDir()
			spillDir := t.TempDir()
			t.Setenv("TMPDIR", spillDir)
			t.Setenv("TMP", spillDir)
			t.Setenv("TEMP", spillDir)
			writeFixture(t, rootPath, "file.txt", "needle short\n"+longLine+"\n"+tc.tail)
			root, err := openDirectory(rootPath)
			if err != nil {
				t.Fatal(err)
			}
			defer root.Close()

			visits := 0
			found, err := searchFiles(root, []string{"file.txt"}, []byte("needle"), func(match searchMatch) error {
				if entries, err := os.ReadDir(spillDir); err != nil || len(entries) != 1 {
					t.Errorf("expected one active spill: %v, error %v", entries, err)
				}
				visits++
				wantText := "needle short"
				if visits == 2 {
					wantText = longLine
				}
				if match.Path != "file.txt" || match.Line != visits || match.Text != wantText {
					t.Errorf("match %d: wrong path, line number or content", visits)
				}
				return tc.writeErr
			})
			if found != tc.found || visits != tc.visits || !errors.Is(err, tc.writeErr) {
				t.Errorf("found %t, visited %d, error %v; want %t, %d, %v", found, visits, err, tc.found, tc.visits, tc.writeErr)
			}
			if entries, err := os.ReadDir(spillDir); err != nil || len(entries) != 0 {
				t.Errorf("spill directory after search: %v, error %v", entries, err)
			}
		})
	}
}
