package parser

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
	Key      *Literal    `json:"key,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

func (n *Node) Node() *Node {
	return n.Value.(*Node)
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
