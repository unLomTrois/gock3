package tokens

import (
	"ck3-parser/internal/app/files"
	"fmt"
	"strconv"
)

type Token struct {
	Value string    `json:"value"`
	Type  TokenType `json:"type"`
	Loc   files.Loc `json:"-"`
}

func New(value string, tokenType TokenType, loc files.Loc) *Token {
	return &Token{
		Value: value,
		Type:  tokenType,
		Loc:   loc,
	}
}

func (t *Token) IsBV() {}

func (t *Token) String() string {
	return fmt.Sprintf("type:\t%v,\tvalue:\t%v", t.Type, strconv.Quote(string(t.Value)))
}
