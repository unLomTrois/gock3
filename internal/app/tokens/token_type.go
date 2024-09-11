package tokens

type TokenType uint8

const (
	COMMENT TokenType = iota
	WORD
	QUOTED_STRING
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

var TokenTypeRegexMap = map[TokenType]string{
	COMMENT:       `^#(.+)?`,
	WORD:          `^@?(?:[\w-]+:)?[\w.-]+`,
	QUOTED_STRING: `^"(.*?)"`,
	NUMBER:        `^-?\d+([.,]\d+)?`,
	BOOL:          `^(yes|no)`,
	NEXTLINE:      `^\n+`,
	EQUALS:        `^==?`,
	START:         `^{`,
	END:           `^}`,
	WHITESPACE:    `^ +`,
	TAB:           `^\t+`,
	COMPARISON:    `^[\<\>]=?`,
}

// TokenCheckOrder defines the order in which tokens should be checked
var TokenCheckOrder = []TokenType{
	WHITESPACE,
	TAB,
	NEXTLINE,
	COMPARISON,
	COMMENT,
	QUOTED_STRING,
	BOOL,
	NUMBER,
	WORD,
	EQUALS,
	START,
	END,
}

func (tt TokenType) String() string {
	switch tt {
	case COMMENT:
		return "COMMENT"
	case WORD:
		return "WORD"
	case QUOTED_STRING:
		return "QUOTED_STRING"
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

func (tt TokenType) MarshalText() ([]byte, error) {
	return []byte(tt.String()), nil
}
