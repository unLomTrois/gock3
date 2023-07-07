package parser2

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
	Value    interface{} `json:"value,omitempty"`
}

func (n *Node) Node() *Node {
	return n.Value.(*Node)
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

	switch t := n.Value.(type) {
	case string:
		return []byte(t)
	case float32:
		return []byte(fmt.Sprintf("%g", t))
	}

	return nil
}
