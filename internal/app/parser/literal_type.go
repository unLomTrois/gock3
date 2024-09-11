package parser

type LiteralType string

const (
	NumberLiteral  LiteralType = "NumberLiteral"
	BoolLiteral    LiteralType = "BoolLiteral"
	StringLiteral  LiteralType = "StringLiteral"
	WordLiteral    LiteralType = "WordLiteral"
	CommentLiteral LiteralType = "CommentLiteral"
)
