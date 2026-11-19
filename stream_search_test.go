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
		data := "needle before error\n" + tail + "data beyond first buffer\n"
		source := io.MultiReader(
			strings.NewReader(data),
			iotest.ErrReader(failure),
		)
		var staged stagedMatches
		t.Cleanup(func() {
			if err := staged.reset(); err != nil {
				t.Error(err)
			}
		})
		valid, err := collectMatches(bufio.NewReaderSize(source, 16), strings.NewReader(data), []byte("needle"), &staged)
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
		source := strings.NewReader(text)
		gotValid, err := collectMatches(bufio.NewReader(source), source, []byte("needle"), &staged)
		if valid != gotValid || (valid && !errors.Is(err, os.ErrClosed)) || (!valid && err != nil) {
			t.Errorf("valid source %t: valid %t, error %v", valid, gotValid, err)
		}
	}
}

func TestLongLineStagingRejectsChangedBytes(t *testing.T) {
	for _, tc := range []struct {
		name string
		size int
	}{
		{name: "memory", size: 2 * searchBufferSize},
		{name: "spill", size: matchMemoryLimit},
	} {
		t.Run(tc.name, func(t *testing.T) {
			line := "needle " + strings.Repeat("x", tc.size)
			for _, suffix := range []string{"a\n", "a\n\x00\n"} {
				probe := bufio.NewReaderSize(strings.NewReader(line+suffix), searchBufferSize)
				readerAt := strings.NewReader(line + "b\n")
				var matches stagedMatches
				valid, err := collectMatches(probe, readerAt, []byte("needle"), &matches)
				if suffix == "a\n" {
					if !valid || err == nil || !strings.Contains(err.Error(), "source changed while staging long line 1") {
						t.Errorf("changed line was eligible for replay: valid %t, error %v", valid, err)
					}
				} else if valid || err != nil {
					t.Errorf("later NUL should skip file despite changed bytes: valid %t, error %v", valid, err)
				}
				if err := matches.reset(); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestLongUnmatchedLineDoesNotStage(t *testing.T) {
	source := strings.NewReader(strings.Repeat("x", 8*searchBufferSize))
	var matches stagedMatches
	// Passing no ReaderAt also proves the unmatched content is never reread.
	valid, err := collectMatches(bufio.NewReaderSize(source, searchBufferSize), nil, []byte("needle"), &matches)
	if !valid || err != nil || matches.spill != nil || cap(matches.memory) != 0 {
		t.Fatalf("unmatched long line: valid %t, error %v, spill %v, stage capacity %d", valid, err, matches.spill, cap(matches.memory))
	}
}

func TestLongLiteralCrossesSeveralFragments(t *testing.T) {
	literal := strings.Repeat("☕abc", 12)
	line := "123456789012345" + literal + " end"
	source := strings.NewReader("header\n" + line + "\r\ntrailer\n")
	reader := bufio.NewReaderSize(source, 16)
	var matches stagedMatches
	if valid, err := collectMatches(reader, source, []byte(literal), &matches); !valid || err != nil || matches.spill != nil {
		t.Fatalf("long literal: valid %t, error %v, spill %v", valid, err, matches.spill)
	}
	found, err := matches.replay("file.txt", reader, func(match searchMatch) error {
		if match.Path != "file.txt" || match.Line != 2 || match.Text != line {
			t.Errorf("wrong match: %#v", match)
		}
		return nil
	})
	if !found || err != nil {
		t.Errorf("replay found %t, error %v", found, err)
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
