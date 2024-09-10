package lexer

import (
	"fmt"
	"regexp"
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

// CompileRegexes compiles the regular expressions from TokenTypeToRegex map
func CompileRegexes() map[TokenType]*regexp.Regexp {
	var CompiledRegexMap = make(map[TokenType]*regexp.Regexp)

	for tokenType, regexStr := range TokenTypeToRegex {
		regex, err := regexp.Compile(regexStr)
		if err != nil {
			panic(fmt.Sprintf("Failed to compile regex for TokenType %s: %s", tokenType, err))
		}
		CompiledRegexMap[tokenType] = regex
	}

	return CompiledRegexMap
}
