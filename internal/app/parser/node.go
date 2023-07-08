package parser

import (
	"fmt"
)

type NodeType string

const (
	NextLine   NodeType = "NextLine"
	Comment    NodeType = "Comment"
	Entity     NodeType = "Entity"
	Block      NodeType = "Block"
	Script     NodeType = "Script"
	Property   NodeType = "Property"
	Comparison NodeType = "Comparison"
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

type Node struct {
	// Parent *any     `json:"-"`
	Type     NodeType    `json:"type"`
	Key      *Literal    `json:"key,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

func (n *Node) Node() *Node {
	return n.Value.(*Node)
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

func (n *Node) KeyLiteral() []byte {
	return []byte(n.Key.String())
}

func (n *Node) DataLiteral() []byte {
	if n.Type == Comment {
		return []byte(n.Value.(string))
	}
	return []byte(n.Value.(*Literal).String())
}
