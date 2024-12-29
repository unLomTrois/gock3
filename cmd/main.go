package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unLomTrois/gock3/internal/app/cli"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if err := root(os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func root(args []string) error {
	commands := []cli.Command{
		cli.NewParseCommand(),
		cli.NewProjectCommand(),
	}

	if len(args) < 2 {
		fmt.Println("No command provided")
		printHelp(commands)
		return nil
	}

	subcommand := args[1]

	actualArgs := args[2:]

	// Run the command
	for _, cmd := range commands {
		if subcommand == cmd.Name() {
			return cmd.Run(actualArgs)
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}

// todo: add descriptions for each command
func printHelp(commands []cli.Command) {
	fmt.Println("Available commands:")

	for _, cmd := range commands {
		fmt.Printf("  %s - %s\n", cmd.Name(), cmd.Description())
	}
}
