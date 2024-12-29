package cli

import (
	"flag"

	"github.com/unLomTrois/gock3/pkg/project"
)

type ProjectCommand struct {
	fs             *flag.FlagSet
	game_dir       string
	mod_descriptor string
}

func NewProjectCommand() *ProjectCommand {
	command := &ProjectCommand{
		fs: flag.NewFlagSet("project", flag.ExitOnError),
	}

	command.fs.StringVar(
		&command.game_dir,
		"game",
		"",
		`gock3 project --game "steamapps/common/Crusader Kings III/game"`,
	)

	command.fs.StringVar(
		&command.mod_descriptor,
		"mod",
		"",
		`gock3 project --mod "Documents/Paradox Interactive/Crusader Kings III/mod/<modname>.mod"`,
	)

	return command
}

func (c *ProjectCommand) Name() string {
	return c.fs.Name()
}

func (c *ProjectCommand) Description() string {
	return "Load the whole project and parses them to data structures"
}

// Run parses the input file and generates the output files
// args is the list of command line arguments
// The first argument is the file path
func (c *ProjectCommand) Run(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	project, err := project.NewProject(c.game_dir, c.mod_descriptor)
	if err != nil {
		return err
	}

	project.Load()

	return nil
}

// func (c *ProjectCommand) parse(fullpath string) error {
// 	file_entry := files.NewFileEntry(fullpath, files.FileKind(files.Mod))

// 	ast, err := pdxfile.ParseFile(file_entry)
// 	if err != nil {
// 		return err
// 	}

// 	if c.astFilepath != "" {
// 		if err := utils.SaveJSON(ast, c.astFilepath); err != nil {
// 			return fmt.Errorf("saving parse tree: %w", err)
// 		}
// 		return err
// 	}

// 	ast_string, err := json.MarshalIndent(ast, "", "  ")
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println(string(ast_string))

// 	return nil
// }
