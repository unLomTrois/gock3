package parser

type Literal struct {
	Type  LiteralType `json:"type"`
	Value interface{} `json:"value"`
}

func (l *Literal) String() string {
	switch l.Type {
	case StringLiteral, CommentLiteral, WordLiteral, BoolLiteral, NumberLiteral:
		return l.Value.(string)
	default:
		panic("Unknown literal type:" + l.Type.String())
	}
}
