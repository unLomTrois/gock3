package parser

import (
	"fmt"
)

type LiteralType string

const (
	NumberLiteral  LiteralType = "NumberLiteral"
	BoolLiteral    LiteralType = "BoolLiteral"
	StringLiteral  LiteralType = "StringLiteral"
	WordLiteral    LiteralType = "WordLiteral"
	CommentLiteral LiteralType = "CommentLiteral"
)

type Literal struct {
	Type  LiteralType `json:"type"`
	Value interface{} `json:"value"`
}

func (l *Literal) String() string {
	switch l.Type {
	case StringLiteral, CommentLiteral, WordLiteral, BoolLiteral:
		return l.Value.(string)
	case NumberLiteral:
		return fmt.Sprintf("%g", l.Value.(float64))
	default:
		panic("Unknown literal type:" + l.Type)
	}
}
