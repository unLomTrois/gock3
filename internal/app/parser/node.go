package parser

import "ck3-parser/internal/app/tokens"

type Node struct {
	Key      *tokens.Token `json:"key,omitempty"`
	Operator *tokens.Token `json:"operator,omitempty"`
	Value    interface{}   `json:"value,omitempty"`
}

func (n *Node) Node() *Node {
	return n.Value.(*Node)
}
