package cli

type Command interface {
	Run(args []string) error

	// Name of the command
	Name() string

	// Description for help
	Description() string
}
