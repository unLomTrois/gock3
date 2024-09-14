package parser

type Node struct {
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
	return []byte(n.Value.(*Literal).String())
}
