package cli

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/pdxfile"
	"github.com/unLomTrois/gock3/internal/app/utils"
)

type ParseCommand struct {
	flagset     *flag.FlagSet
	astFilepath string
}

func NewParseCommand() *ParseCommand {
	command := &ParseCommand{
		flagset: flag.NewFlagSet("parse", flag.ExitOnError),
	}

	command.flagset.StringVar(
		&command.astFilepath,
		"save-ast",
		"",
		"Save the AST to a file\ngock3 parse file.txt --save-ast ast.json",
	)

	return command
}

func (command *ParseCommand) Name() string {
	return command.flagset.Name()
}

func (command *ParseCommand) Description() string {
	return "Parse a file and generate the output files"
}

// Run parses the input file and generates the output files
// args is the list of command line arguments
// The first argument is the file path
func (command *ParseCommand) Run(args []string) error {
	if err := command.validateArgs(args); err != nil {
		return err
	}

	if err := command.flagset.Parse(args[1:]); err != nil {
		return err
	}

	filePath := args[0]
	fullpath, err := utils.FileExists(filePath)
	if err != nil {
		return err
	}

	return command.parse(fullpath)
}

func (command *ParseCommand) validateArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("not enough arguments")
	}
	return nil
}

func (command *ParseCommand) parse(fullpath string) error {
	fileEntry := files.NewFileEntry(fullpath, files.FileKind(files.Mod))

	ast, err := pdxfile.ParseFile(fileEntry)
	if err != nil {
		return err
	}

	if command.astFilepath == "" {
		return nil
	}

	if err := utils.SaveJSON(ast, command.astFilepath); err != nil {
		return fmt.Errorf("failed to save AST: %w", err)
	}

	absPath, err := filepath.Abs(command.astFilepath)
	if err != nil {
		return err
	}

	log.Println("Saved parse tree to", absPath)
	return nil
}
