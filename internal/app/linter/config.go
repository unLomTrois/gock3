package linter

type IntendStyle string

const (
	TABS   IntendStyle = "tabs"
	SPACES IntendStyle = "spaces"
)

type LintConfig struct {
	IntendStyle            IntendStyle
	IntendSize             int
	TrimTrailingWhitespace bool
	InsertFinalNewline     bool
	CharSet                string
	EndOfLine              []byte
}
