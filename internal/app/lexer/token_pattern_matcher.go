package lexer

import (
	"fmt"
	"log"
	"regexp"
)

// TokenPatternMatcher is a structure for working with token regular expressions
type TokenPatternMatcher struct {
	compiledRegexMap map[TokenType]*regexp.Regexp
}

// NewTokenPatternMatcher creates a new instance of TokenPatternMatcher and compiles regular expressions
func NewTokenPatternMatcher() *TokenPatternMatcher {
	tpm := &TokenPatternMatcher{
		compiledRegexMap: make(map[TokenType]*regexp.Regexp),
	}
	tpm.compileRegexes()
	return tpm
}

// compileRegexes compiles regular expressions and stores them in the map
func (tpm *TokenPatternMatcher) compileRegexes() {
	tokenTypeRegexMap := map[TokenType]string{
		COMMENT:    `^#(.+)?`,
		WORD:       `^(?:\w+:)?\w+(?:\.\w+)*`,
		STRING:     `^"(.*?)"`,
		NUMBER:     `^-?\d+([.,]\d+)?`,
		BOOL:       `^(yes|no)`,
		NEXTLINE:   `^\n+`,
		EQUALS:     `^==?`,
		START:      `^{`,
		END:        `^}`,
		WHITESPACE: `^ +`,
		TAB:        `^\t+`,
		COMPARISON: `^[\<\>]=?`,
	}

	for tokenType, regexPattern := range tokenTypeRegexMap {
		regex, err := regexp.Compile(regexPattern)
		if err != nil {
			log.Fatalf("Failed to compile regex for TokenType %s: %v", tokenType, err)
		}
		tpm.compiledRegexMap[tokenType] = regex
	}
}

// MatchToken finds the first match for the given token type and text
func (tpm *TokenPatternMatcher) MatchToken(tokenType TokenType, text []byte) []byte {
	regex, exists := tpm.compiledRegexMap[tokenType]
	if !exists {
		fmt.Printf("No regex found for TokenType %s\n", tokenType)
		return nil
	}
	return regex.Find(text)
}
