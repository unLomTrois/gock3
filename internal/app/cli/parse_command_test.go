// parse_command_test.go
package cli_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/unLomTrois/gock3/internal/app/cli"
)

// --------------------
// Unit-test style tests
// --------------------

func TestNewParseCommand(t *testing.T) {
	tests := []struct {
		name string
		want *cli.ParseCommand
	}{
		{
			name: "Create new parse command",
			want: &cli.ParseCommand{}, // We can’t compare the FlagSet directly, so we mainly check type here
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cli.NewParseCommand()
			// We can’t do a direct DeepEqual on *flag.FlagSet,
			// so we do a simpler check on type for demonstration:
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewParseCommand() = %T, want %T", got, tt.want)
			}
		})
	}
}

func TestParseCommand_Name(t *testing.T) {
	cmd := cli.NewParseCommand()
	expectedName := "parse"
	if cmd.Name() != expectedName {
		t.Errorf("ParseCommand.Name() = %v, want %v", cmd.Name(), expectedName)
	}
}

func TestParseCommand_Description(t *testing.T) {
	cmd := cli.NewParseCommand()
	expectedDesc := "Parse a file and generate the output files"
	if cmd.Description() != expectedDesc {
		t.Errorf("ParseCommand.Description() = %v, want %v", cmd.Description(), expectedDesc)
	}
}

// --------------------------
// Integration-style tests
// --------------------------

func TestParseCommand_MissingArguments(t *testing.T) {
	// If no arguments are passed, we expect an error because
	// the command does: file_path := args[0]
	cmd := cli.NewParseCommand()
	err := cmd.Run([]string{}) // no args
	if err == nil {
		t.Errorf("expected error for missing arguments, got nil")
	}
}

func TestParseCommand_NonExistingFile(t *testing.T) {
	// If the file doesn’t exist, FileExists should return an error.
	cmd := cli.NewParseCommand()
	err := cmd.Run([]string{"this_file_does_not_exist.txt"})
	if err == nil {
		t.Errorf("expected error for non-existing file, got nil")
	}
}

func TestParseCommand_ValidFile_NoAstFlag(t *testing.T) {
	// Create a temporary file to parse
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up

	// Write something minimal to the file so pdxfile.ParseFile won't choke
	_, err = tmpFile.WriteString("hello = world")
	if err != nil {
		t.Fatal(err)
	}
	// Close before we parse it
	tmpFile.Close()

	cmd := cli.NewParseCommand()
	// The first element in args is the file, subsequent elements are flags
	err = cmd.Run([]string{tmpFile.Name()})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestParseCommand_ValidFile_WithAstFlag(t *testing.T) {
	// Create a temporary file to parse
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up

	_, err = tmpFile.WriteString("hello = world")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Create a path where we want to save our AST
	tmpDir := t.TempDir()
	astPath := filepath.Join(tmpDir, "ast.json")

	cmd := cli.NewParseCommand()
	// Provide the file and the --save-ast flag
	err = cmd.Run([]string{tmpFile.Name(), "--save-ast", astPath})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify that the AST file was created
	if _, statErr := os.Stat(astPath); os.IsNotExist(statErr) {
		t.Errorf("expected AST file to exist at %s, but it does not", astPath)
	}
}
