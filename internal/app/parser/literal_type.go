package parser

type LiteralType uint8

const (
	NumberLiteral LiteralType = iota
	BoolLiteral
	StringLiteral
	WordLiteral
	CommentLiteral
)

func (lt LiteralType) String() string {
	switch lt {
	case NumberLiteral:
		return "NumberLiteral"
	case BoolLiteral:
		return "BoolLiteral"
	case StringLiteral:
		return "StringLiteral"
	case WordLiteral:
		return "WordLiteral"
	case CommentLiteral:
		return "CommentLiteral"
	default:
		return "UnknownLiteral"
	}
}

func (lt LiteralType) MarshalText() ([]byte, error) {
	return []byte(lt.String()), nil
}
