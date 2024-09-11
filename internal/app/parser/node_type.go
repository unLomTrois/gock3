package parser

type NodeType uint8

const (
	NextLine NodeType = iota
	Comment
	Entity
	Block
	Script
	Property
	Comparison
)

func (nt NodeType) String() string {
	switch nt {
	case NextLine:
		return "NextLine"
	case Comment:
		return "Comment"
	case Entity:
		return "Entity"
	case Block:
		return "Block"
	case Script:
		return "Script"
	case Property:
		return "Property"
	case Comparison:
		return "Comparison"
	default:
		return "Unknown"
	}
}

func (nt NodeType) MarshalText() ([]byte, error) {
	return []byte(nt.String()), nil
}
