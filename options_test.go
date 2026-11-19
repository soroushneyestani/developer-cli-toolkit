package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestInspectAndOptionsCLI(t *testing.T) {
	binary := buildCLI(t)

	t.Run("search extension filters and JSON", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, ".hidden.go", "needle \"quoted\"\t<tag>\n")
		writeFixture(t, root, "dir name/file.txt", "first\r\nneedle café\r\n")
		writeFixture(t, root, "uppercase.GO", "needle excluded\n")
		writeFixture(t, root, "suffix.go.bak", "needle excluded\n")
		writeFixture(t, root, "archive.tar.gz", "needle compound suffix\n")

		result := runCLI(t, binary, root, "search", "needle", root,
			"--ext", ".go", "--ext", ".txt", "--ext", ".go", "--json")
		assertCLIJSON(t, result,
			map[string]any{"path": ".hidden.go", "line": float64(1), "text": "needle \"quoted\"\t<tag>"},
			map[string]any{"path": "dir name/file.txt", "line": float64(2), "text": "needle café"},
		)

		result = runCLI(t, binary, root, "search", "needle", "--ext", ".tar.gz")
		assertCLI(t, result, 0, "archive.tar.gz:1:needle compound suffix\n", false)
		result = runCLI(t, binary, root, "search", "needle", "--ext", ".md", "--json")
		assertCLI(t, result, 1, "", false)
	})

	t.Run("malformed options on either command", func(t *testing.T) {
		root := t.TempDir()
		for _, filter := range []string{"go", "", ".go/sub", `.go\sub`} {
			t.Run(fmt.Sprintf("extension %q", filter), func(t *testing.T) {
				result := runCLI(t, binary, root, "search", "needle", "--ext", filter)
				assertCLI(t, result, 2, "", true)
				result = runCLI(t, binary, root, "inspect", "--ext", filter)
				assertCLI(t, result, 2, "", true)
			})
		}
		result := runCLI(t, binary, root, "inspect", "--ext")
		assertCLI(t, result, 2, "", true)
	})

	t.Run("inspect relative names and byte sizes", func(t *testing.T) {
		root := t.TempDir()
		const hidden = "\x00\xff"
		const nested = "café\r\n"
		writeFixture(t, root, ".hidden.dat", hidden)
		writeFixture(t, root, "dir name/nested.txt", nested)
		writeFixture(t, root, "dir name/zero.txt", "")
		if err := os.Mkdir(filepath.Join(root, "empty"), 0o700); err != nil {
			t.Fatalf("create empty directory: %v", err)
		}

		result := runCLI(t, binary, root, "inspect")
		assertCLI(t, result, 0, fmt.Sprintf(".hidden.dat\t%d\ndir name/nested.txt\t%d\ndir name/zero.txt\t0\n", len(hidden), len(nested)), false)

		result = runCLI(t, binary, root, "inspect", root, "--ext", ".txt", "--json")
		assertCLIJSON(t, result,
			map[string]any{"path": "dir name/nested.txt", "bytes": float64(len(nested))},
			map[string]any{"path": "dir name/zero.txt", "bytes": float64(0)},
		)
		result = runCLI(t, binary, root, "inspect", filepath.Join(root, "empty"), "--json")
		assertCLI(t, result, 0, "", false)
	})

	t.Run("inspect JSON escapes filename delimiters", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Windows filenames cannot contain colons or tabs")
		}
		root := t.TempDir()
		const name = "colon:tab\t.txt"
		writeFixture(t, root, name, "x")
		result := runCLI(t, binary, root, "inspect", "--json")
		assertCLIJSON(t, result, map[string]any{"path": name, "bytes": float64(1)})
	})

	t.Run("inspect excludes links", func(t *testing.T) {
		root := t.TempDir()
		file := writeFixture(t, root, "real.txt", "abc")
		if err := os.Symlink(file, filepath.Join(root, "linked.txt")); err != nil {
			t.Skipf("cannot create file symlink on this host: %v", err)
		}
		result := runCLI(t, binary, root, "inspect")
		assertCLI(t, result, 0, "real.txt\t3\n", false)
	})

	t.Run("inspect misuse", func(t *testing.T) {
		root := t.TempDir()
		file := writeFixture(t, root, "file.txt", "x")
		cases := []struct {
			name string
			args []string
		}{
			{"unknown flag", []string{"inspect", "--unrecognized"}},
			{"extra argument", []string{"inspect", root, "extra"}},
			{"missing root", []string{"inspect", filepath.Join(root, "missing")}},
			{"file root", []string{"inspect", file}},
		}
		for _, test := range cases {
			t.Run(test.name, func(t *testing.T) {
				result := runCLI(t, binary, root, test.args...)
				assertCLI(t, result, 2, "", true)
			})
		}
	})
}

func assertCLIJSON(t *testing.T, result cliResult, want ...map[string]any) {
	t.Helper()
	if result.exit != 0 || result.stderr != "" {
		t.Errorf("exit %d, stderr %q; want exit 0 and empty stderr", result.exit, result.stderr)
	}
	output, terminated := strings.CutSuffix(result.stdout, "\n")
	if !terminated {
		t.Fatalf("JSON output lacks its final newline: %q", result.stdout)
	}
	lines := strings.Split(output, "\n")
	if len(lines) != len(want) {
		t.Fatalf("got %d JSON records, want %d (stdout %q)", len(lines), len(want), result.stdout)
	}
	for i, line := range lines {
		var got map[string]any
		if err := json.Unmarshal([]byte(line), &got); err != nil {
			t.Errorf("JSON record %d: %v", i+1, err)
			continue
		}
		if !reflect.DeepEqual(got, want[i]) {
			t.Errorf("JSON record %d: got %#v, want %#v", i+1, got, want[i])
		}
	}
}
