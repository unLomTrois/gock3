package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/pdxfile"
	"github.com/unLomTrois/gock3/internal/app/utils"
)

type ParseCommand struct {
	fs          *flag.FlagSet
	astFilepath string
}

func NewParseCommand() *ParseCommand {
	command := &ParseCommand{
		fs: flag.NewFlagSet("parse", flag.ExitOnError),
	}

	command.fs.StringVar(
		&command.astFilepath,
		"save-ast",
		"",
		"Save the AST to a file\ngock3 parse file.txt --save-ast ast.json",
	)

	return command
}

func (c *ParseCommand) Name() string {
	return c.fs.Name()
}

// Run parses the input file and generates the output files
// args is the list of command line arguments
// The first argument is the file path
func (c *ParseCommand) Run(args []string) error {
	if err := c.fs.Parse(args[1:]); err != nil {
		return err
	}

	file_path := args[0]
	fullpath, err := utils.FileExists(file_path)
	if err != nil {
		return err
	}

	return c.parse(fullpath)
}

func (c *ParseCommand) parse(fullpath string) error {
	file_entry := files.NewFileEntry(fullpath, files.FileKind(files.Mod))

	ast, err := pdxfile.ParseFile(file_entry)
	if err != nil {
		return err
	}

	if c.astFilepath != "" {
		if err := utils.SaveJSON(ast, c.astFilepath); err != nil {
			return fmt.Errorf("saving parse tree: %w", err)
		}
		return err
	}

	ast_string, err := json.MarshalIndent(ast, "", "  ")
	if err != nil {
		return err
	}

	log.Println(string(ast_string)[0])

	return nil
}
