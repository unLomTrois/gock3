package lexer

import (
	"fmt"
	"strconv"
)

type Token struct {
	Type  TokenType `json:"type"`
	Value string    `json:"value"`
}

func (t *Token) String() string {
	return fmt.Sprintf("type:\t%v,\tvalue:\t%v", t.Type, strconv.Quote(string(t.Value)))
}
