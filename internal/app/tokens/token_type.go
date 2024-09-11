package tokens

type TokenType string

// Grouping constants for better readability
const (
	COMMENT    TokenType = "COMMENT"
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
var TokenCheckOrder = []TokenType{
	WHITESPACE,
	TAB,
	NEXTLINE,
	COMPARISON,
	COMMENT,
	STRING,
	BOOL,
	NUMBER,
	WORD,
	EQUALS,
	START,
	END,
}
