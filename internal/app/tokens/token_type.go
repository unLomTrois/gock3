package tokens

type TokenType int

const (
	COMMENT TokenType = iota
	WORD
	STRING
	NUMBER
	BOOL
	NEXTLINE
	EQUALS
	START
	END
	WHITESPACE
	TAB
	COMPARISON
)

func (tt TokenType) String() string {
	switch tt {
	case COMMENT:
		return "COMMENT"
	case WORD:
		return "WORD"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case BOOL:
		return "BOOL"
	case NEXTLINE:
		return "NEXTLINE"
	case EQUALS:
		return "EQUALS"
	case START:
		return "START"
	case END:
		return "END"
	case WHITESPACE:
		return "WHITESPACE"
	case TAB:
		return "TAB"
	case COMPARISON:
		return "COMPARISON"
	default:
		return "UNKNOWN"
	}
}

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
