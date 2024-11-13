package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unLomTrois/gock3/internal/app/cli"
)

func main() {
	if err := root(os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func root(args []string) error {
	if len(args) < 2 {
		fmt.Println("No command provided")
		printHelp()
		return nil
	}

	commands := []cli.Command{
		cli.NewParseCommand(),
		cli.NewProjectCommand(),
	}

	subcommand := args[1]

	// check if subcommand is valid
	for _, cmd := range commands {
		if subcommand == cmd.Name() {
			return cmd.Run(args[2:])
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  parse   - Description of parse command")
	fmt.Println("  project - Description of project command")
	// add descriptions for each command
}

// func run() error {
// 	start := time.Now()
// 	defer func() {
// 		log.Printf("Total execution time: %s", time.Since(start))
// 	}()

// 	vanilla_root := ""
// 	mod_root := ""
// 	replace_paths := []string{}

// 	mod_loader := files.NewModLoader(mod_root, replace_paths)

// 	fset := files.NewFileSet(vanilla_root, mod_loader)

// 	traits_dir := "C:/Users/vadim/Documents/Paradox Interactive/Crusader Kings III/mod/T4N-CK3/T4N/common/traits"

// 	err := fset.Scan(traits_dir)
// 	if err != nil {
// 		return fmt.Errorf("scanning files: %w", err)
// 	}

// 	path := inputFilePath
// 	fullpath, err := filepath.Abs(path)
// 	if err != nil {
// 		return fmt.Errorf("getting absolute path: %w", err)
// 	}
// 	file_entry := files.NewFileEntry(fullpath, files.FileKind(files.Mod))

// 	parseTrees, err := pdxfile.ParseFile(file_entry)
// 	if err != nil {
// 		return fmt.Errorf("parsing tokens: %w", err)
// 	}

// 	if err := os.MkdirAll(outputDir, 0755); err != nil {
// 		return err
// 	}

// 	if err := utils.SaveJSON(parseTrees, filepath.Join(outputDir, parseTreeFile)); err != nil {
// 		return fmt.Errorf("saving parse tree: %w", err)
// 	}

// 	return nil
// }
