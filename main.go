package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

const version = "0.1.1-rc.1"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	exitCode := 0
	var extensions []string
	var asJSON bool
	root := &cobra.Command{
		Use:           "dct",
		Short:         "Search text and inspect workspace files",
		Version:       version,
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE:          func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}
	root.CompletionOptions.DisableDefaultCmd = true
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)

	search := &cobra.Command{
		Use:   "search <literal> [root]",
		Short: "Find a case-sensitive literal in UTF-8 text files",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			if args[0] == "" || !utf8.ValidString(args[0]) {
				return errors.New("search requires a nonempty UTF-8 literal")
			}
			workspace, err := openWorkspace(args[1:], extensions)
			if err != nil {
				return err
			}
			defer workspace.Close()
			files, err := listFiles(workspace, extensions)
			if err != nil {
				return err
			}
			encoder := json.NewEncoder(stdout)
			found, err := searchFiles(workspace, files, []byte(args[0]), func(match searchMatch) error {
				if asJSON {
					return encoder.Encode(match)
				}
				_, err := fmt.Fprintf(stdout, "%s:%d:%s\n", match.Path, match.Line, match.Text)
				return err
			})
			if !found {
				exitCode = 1
			}
			return err
		},
	}

	inspect := &cobra.Command{
		Use:   "inspect [root]",
		Short: "List regular workspace files and their byte sizes",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			workspace, err := openWorkspace(args, extensions)
			if err != nil {
				return err
			}
			defer workspace.Close()
			files, err := listFiles(workspace, extensions)
			if err != nil {
				return err
			}
			encoder := json.NewEncoder(stdout)
			return inspectFiles(workspace, files, func(file fileSize) error {
				if asJSON {
					return encoder.Encode(file)
				}
				_, err := fmt.Fprintf(stdout, "%s\t%d\n", file.Path, file.Bytes)
				return err
			})
		},
	}

	for _, command := range []*cobra.Command{search, inspect} {
		command.Flags().StringArrayVar(&extensions, "ext", nil, "include files ending with this suffix (repeatable)")
		command.Flags().BoolVar(&asJSON, "json", false, "emit newline-delimited JSON")
		root.AddCommand(command)
	}
	if err := root.Execute(); err != nil {
		fmt.Fprintf(stderr, "dct: %v\n", err)
		return 2
	}
	return exitCode
}

func openWorkspace(roots, extensions []string) (*os.Root, error) {
	for _, extension := range extensions {
		if !strings.HasPrefix(extension, ".") || strings.ContainsAny(extension, `/\`) {
			return nil, fmt.Errorf("invalid extension %q: expected a dot-prefixed suffix without path separators", extension)
		}
	}
	root := "."
	if len(roots) > 0 {
		root = roots[0]
		if root == "" {
			return nil, errors.New("root must be a directory path")
		}
	}
	return openDirectory(root)
}
