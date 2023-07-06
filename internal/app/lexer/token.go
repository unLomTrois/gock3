package lexer

import (
	"fmt"
	"strconv"
)

type TokenType string

const (
	COMMENT    TokenType = "COMMENT"
	SCRIPT     TokenType = "SCRIPT"
	WORD       TokenType = "WORD"
	NUMBER     TokenType = "NUMBER"
	NULL       TokenType = "NULL"
	WHITESPACE TokenType = "WHITESPACE"
	NEXTLINE   TokenType = "NEXTLINE"
	TAB        TokenType = "TAB"
	EQUALS     TokenType = "EQUALS"
	START      TokenType = "START"
	END        TokenType = "END"
	COMPARISON TokenType = "COMPARISON"
)

var Spec = map[string]TokenType{
	`^[\<\>]=?`:                  COMPARISON,
	`^#(.+)?`:                    COMMENT,
	`^scripted_(trigger|effect)`: SCRIPT,
	`^(\w+):?[a-zA-Z0-9_]+(\.[a-zA-Z0-9_]+)*`: WORD,
	`^"(.*?)"`:           WORD,
	`^-?\d+[\.,]?(\d?)+`: NUMBER,
	`^ +`:                NULL,
	`^\n+`:               NEXTLINE,
	`^\t+`:               NULL,
	`^==?`:               EQUALS,
	`^{`:                 START,
	`^}`:                 END,
}

type Token struct {
	Type  TokenType `json:"type"`
	Value []byte    `json:"value"`
}

func (t Token) String() string {
	return fmt.Sprintf("type:\t%v,\tvalue:\t%v", t.Type, strconv.Quote(string(t.Value)))
}
