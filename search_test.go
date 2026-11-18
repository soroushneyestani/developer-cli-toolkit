package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestSearchCLI(t *testing.T) {
	binary := buildCLI(t)

	t.Run("literal records in path order", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, "a.txt", "a+b and a+b\r\nA+B\ncafé a+b\r\r\nunterminated a+b")
		writeFixture(t, root, "folder/inside.txt", "a+b inside\n")
		writeFixture(t, root, "folder-other.txt", "a+b outside\n")
		writeFixture(t, root, ".hidden", "a+b hidden\n")
		writeFixture(t, root, "regex.txt", "aaab\n")
		writeFixture(t, root, ".gitignore", "folder/\n")

		result := runCLI(t, binary, root, "search", "a+b")
		assertCLI(t, result, 0, ".hidden:1:a+b hidden\n"+
			"a.txt:1:a+b and a+b\n"+
			"a.txt:3:café a+b\r\n"+
			"a.txt:4:unterminated a+b\n"+
			"folder-other.txt:1:a+b outside\n"+
			"folder/inside.txt:1:a+b inside\n", false)

		result = runCLI(t, binary, root, "search", "café")
		assertCLI(t, result, 0, "a.txt:3:café a+b\r\n", false)
	})

	t.Run("invalid files emit no prior matches", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, "a-nul.txt", "needle before NUL\n\x00\n")
		writeFixture(t, root, "b-utf8.txt", "needle before invalid UTF-8\n\xff\n")
		writeFixture(t, root, "c-valid.txt", "needle in valid file")

		result := runCLI(t, binary, root, "search", "needle", root)
		assertCLI(t, result, 0, "c-valid.txt:1:needle in valid file\n", false)
	})

	t.Run("long lines", func(t *testing.T) {
		root := t.TempDir()
		line := strings.Repeat("x", 70*1024) + "a needle"
		writeFixture(t, root, "long.txt", line+"\n")

		result := runCLI(t, binary, root, "search", "needle")
		assertCLI(t, result, 0, "long.txt:1:"+line+"\n", false)
	})

	t.Run("unmatched search", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, "file.txt", "unrelated\n")

		result := runCLI(t, binary, root, "search", "needle", root)
		assertCLI(t, result, 1, "", false)
	})

	t.Run("misuse", func(t *testing.T) {
		root := t.TempDir()
		file := writeFixture(t, root, "root.txt", "needle\n")
		cases := []struct {
			name string
			args []string
		}{
			{"missing literal", []string{"search"}},
			{"empty literal", []string{"search", ""}},
			{"extra argument", []string{"search", "needle", root, "extra"}},
			{"unknown command", []string{"unknown-command"}},
			{"unknown flag", []string{"search", "needle", "--unrecognized"}},
			{"missing root", []string{"search", "needle", filepath.Join(root, "missing")}},
			{"file root", []string{"search", "needle", file}},
		}
		for _, test := range cases {
			t.Run(test.name, func(t *testing.T) {
				result := runCLI(t, binary, root, test.args...)
				assertCLI(t, result, 2, "", true)
			})
		}
	})

	t.Run("help and version", func(t *testing.T) {
		root := t.TempDir()
		for _, args := range [][]string{
			{"--help"},
			{"search", "needle", filepath.Join(root, "missing"), "--help"},
			{"--version"},
		} {
			result := runCLI(t, binary, root, args...)
			if result.exit != 0 || strings.TrimSpace(result.stdout) == "" || result.stderr != "" {
				t.Errorf("%q: exit %d, stdout %q, stderr %q; want success with output only on stdout", args, result.exit, result.stdout, result.stderr)
			}
		}
	})

	t.Run("symlink traversal and root rejection", func(t *testing.T) {
		root := t.TempDir()
		targetRoot := t.TempDir()
		targetFile := writeFixture(t, targetRoot, "target.txt", "needle outside root\n")
		if err := os.Symlink(targetFile, filepath.Join(root, "file-link.txt")); err != nil {
			t.Skipf("cannot create file symlink on this host: %v", err)
		}
		dirLink := filepath.Join(root, "dir-link")
		if err := os.Symlink(targetRoot, dirLink); err != nil {
			t.Skipf("cannot create directory symlink on this host: %v", err)
		}

		result := runCLI(t, binary, root, "search", "needle")
		assertCLI(t, result, 1, "", false)

		result = runCLI(t, binary, root, "search", "needle", dirLink)
		assertCLI(t, result, 2, "", true)
	})

	t.Run("unreadable file is an error", func(t *testing.T) {
		root := t.TempDir()
		file := writeFixture(t, root, "unreadable.txt", "needle\n")
		if err := os.Chmod(file, 0); err != nil {
			t.Skipf("cannot change file permissions on this host: %v", err)
		}
		t.Cleanup(func() {
			if err := os.Chmod(file, 0o600); err != nil {
				t.Errorf("restore fixture permissions: %v", err)
			}
		})
		probe, err := os.Open(file)
		if err == nil {
			if closeErr := probe.Close(); closeErr != nil {
				t.Fatalf("close permission probe: %v", closeErr)
			}
			t.Skip("this host still permits reading the file after chmod 000")
		}
		if !errors.Is(err, os.ErrPermission) {
			t.Fatalf("probing unreadable file: %v", err)
		}

		result := runCLI(t, binary, root, "search", "needle")
		assertCLI(t, result, 2, "", true)
	})
}

type cliResult struct {
	exit   int
	stdout string
	stderr string
}

func buildCLI(t *testing.T) string {
	t.Helper()
	name := "dct"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binary := filepath.Join(t.TempDir(), name)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-mod=readonly", "-buildvcs=false", "-o", binary, ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}
	return binary
}

func runCLI(t *testing.T, binary, workingDirectory string, args ...string) cliResult {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Dir = workingDirectory
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	result := cliResult{}
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if ctx.Err() != nil || !errors.As(err, &exitErr) {
			t.Fatalf("run CLI %q: %v", args, err)
		}
		result.exit = exitErr.ExitCode()
	}
	result.stdout = stdout.String()
	result.stderr = stderr.String()
	return result
}

func writeFixture(t *testing.T, root, relativePath, content string) string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture %q: %v", relativePath, err)
	}
	return path
}

func assertCLI(t *testing.T, result cliResult, wantExit int, wantStdout string, wantDiagnostic bool) {
	t.Helper()
	if result.exit != wantExit {
		t.Errorf("exit %d, want %d (stderr %q)", result.exit, wantExit, result.stderr)
	}
	if result.stdout != wantStdout {
		t.Errorf("stdout (%d bytes) %.200q, want (%d bytes) %.200q", len(result.stdout), result.stdout, len(wantStdout), wantStdout)
	}
	if wantDiagnostic && strings.TrimSpace(result.stderr) == "" {
		t.Error("expected a diagnostic on stderr")
	}
	if !wantDiagnostic && result.stderr != "" {
		t.Errorf("stderr %q, want empty", result.stderr)
	}
}
