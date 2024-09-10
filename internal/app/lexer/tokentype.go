package lexer

type TokenType string

// Grouping constants for better readability
const (
	COMMENT    TokenType = "COMMENT"
	SCRIPT     TokenType = "SCRIPT"
	WORD       TokenType = "WORD"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"
	BOOL       TokenType = "BOOL"
	NEXTLINE   TokenType = "NEXTLINE"
	EQUALS     TokenType = "EQUALS"
	START      TokenType = "START"
	END        TokenType = "END"
	WHITESPACE TokenType = "WHITESPACE"
	TAB        TokenType = "TAB"
	COMPARISON TokenType = "COMPARISON"
)

// TokenCheckOrder defines the order in which tokens should be checked
var tokenCheckOrder = []TokenType{
	WHITESPACE,
	TAB,
	NEXTLINE,
	COMPARISON,
	COMMENT,
	SCRIPT,
	STRING,
	BOOL,
	NUMBER,
	WORD,
	EQUALS,
	START,
	END,
}
