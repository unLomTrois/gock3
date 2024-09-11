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
