package parser

import "fmt"

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

type Node struct {
	// Parent *any     `json:"-"`
	Type     NodeType    `json:"type"`
	Key      interface{} `json:"key,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Data     interface{} `json:"value,omitempty"`
}

func (n *Node) Node() *Node {
	return n.Data.(*Node)
}

func (n *Node) KeyLiteral() []byte {
	switch t := n.Key.(type) {
	case string:
		return []byte(t)
	case float32:
		return []byte(fmt.Sprintf("%g", t))
	}
	return nil
}

func (n *Node) DataLiteral() []byte {

	switch t := n.Data.(type) {
	case string:
		return []byte(t)
	case float32:
		return []byte(fmt.Sprintf("%g", t))
	}

	return nil
}
