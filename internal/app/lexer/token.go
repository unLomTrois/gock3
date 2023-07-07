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
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"
	BOOL       TokenType = "BOOL"
	NEXTLINE   TokenType = "NEXTLINE"
	EQUALS     TokenType = "EQUALS"
	START      TokenType = "START"
	END        TokenType = "END"
	WHITESPACE TokenType = "WHITESPACE"
	TAB        TokenType = "TAB"
	COMPARISON TokenType = "COMPARISON"
)

var TokenTypeToRegex = map[TokenType]string{
	COMMENT:    `^#(.+)?`,
	SCRIPT:     `^scripted_(trigger|effect)`,
	WORD:       `^(?:\w+:)?\w+(?:\.\w+)*`,
	STRING:     `^"(.*?)"`,
	NUMBER:     `^-?\d+[\.,]?(\d?)+`,
	BOOL:       `^(yes|no)`,
	NEXTLINE:   `^\n+`,
	EQUALS:     `^==?`,
	START:      `^{`,
	END:        `^}`,
	WHITESPACE: `^ +`,
	TAB:        `^\t+`,
	COMPARISON: `^[\<\>]=?`,
}

var TokenCheckOrder = []TokenType{
	WHITESPACE,
	TAB,
	NEXTLINE,
	COMPARISON,
	COMMENT,
	SCRIPT,
	STRING,
	BOOL,
	NUMBER,
	WORD,
	EQUALS,
	START,
	END,
}

type Token struct {
	Type  TokenType `json:"type"`
	Value string    `json:"value"`
}

func (t *Token) String() string {
	return fmt.Sprintf("type:\t%v,\tvalue:\t%v", t.Type, strconv.Quote(string(t.Value)))
}
