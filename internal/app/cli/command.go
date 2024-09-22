package cli

type Command interface {
	Run(args []string) error
	Name() string
}
