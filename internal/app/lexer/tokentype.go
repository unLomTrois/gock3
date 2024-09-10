package lexer

import (
	"log"
	"regexp"
)

type TokenType string

// Grouping constants for better readability
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

// Mapping TokenType to respective regex patterns
var tokenTypeRegexMap = map[TokenType]string{
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

// TokenCheckOrder defines the order in which tokens should be checked
var tokenCheckOrder = []TokenType{
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

// CompileRegexes compiles the regular expressions from tokenTypeRegexMap
func CompileRegexes() map[TokenType]*regexp.Regexp {
	compiledRegexMap := make(map[TokenType]*regexp.Regexp)

	for tokenType, regexPattern := range tokenTypeRegexMap {
		regex, err := regexp.Compile(regexPattern)
		if err != nil {
			log.Fatalf("Failed to compile regex for TokenType %s: %v", tokenType, err)
		}
		compiledRegexMap[tokenType] = regex
	}

	return compiledRegexMap
}
