package lexer

import (
	"fmt"
	"regexp"
	"strconv"
)

type Token struct {
	Type  TokenType `json:"type"`
	Value string    `json:"value"`
}

func (t *Token) String() string {
	return fmt.Sprintf("type:\t%v,\tvalue:\t%v", t.Type, strconv.Quote(string(t.Value)))
}

// matchToken finds the first match of the given regular expression in the provided text
// and returns the matched text as a byte slice. If no match is found, it returns an empty byte slice.
func matchToken(reg *regexp.Regexp, text []byte) []byte {
	return reg.Find(text)
}
