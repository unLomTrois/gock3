package tokens

import (
	"fmt"
	"strconv"
)

type Token struct {
	Value string    `json:"value"`
	Type  TokenType `json:"type"`
	Loc   Loc       `json:"-"`
}

func New(value string, tokenType TokenType, loc Loc) *Token {
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

// check if input type is equal to token type
func (t *Token) IsType(input TokenType) bool {
	return t.Type == input
}

// check if token value is value
func (t *Token) Is(value string) bool {
	return t.Value == value
}

// float value
func (t *Token) FloatValue() (float64, error) {
	return strconv.ParseFloat(t.Value, 64)
}
