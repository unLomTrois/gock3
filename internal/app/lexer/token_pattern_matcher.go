package lexer

import (
	"fmt"
	"log"
	"regexp"

	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
)

// TokenPatternMatcher is a structure for working with token regular expressions
type TokenPatternMatcher struct {
	compiledRegexMap map[tokens.TokenType]*regexp.Regexp
}

// NewTokenPatternMatcher creates a new instance of TokenPatternMatcher and compiles regular expressions
func NewTokenPatternMatcher() *TokenPatternMatcher {
	tpm := &TokenPatternMatcher{
		compiledRegexMap: make(map[tokens.TokenType]*regexp.Regexp),
	}
	tpm.compileRegexes()
	return tpm
}

// compileRegexes compiles regular expressions and stores them in the map
func (tpm *TokenPatternMatcher) compileRegexes() {
	for tokenType, regexPattern := range tokens.TokenTypeRegexMap {
		regex, err := regexp.Compile(regexPattern)
		if err != nil {
			log.Fatalf("Failed to compile regex for TokenType %s: %v", tokenType, err)
		}
		tpm.compiledRegexMap[tokenType] = regex
	}
}

// MatchToken finds the first match for the given token type and text
func (tpm *TokenPatternMatcher) MatchToken(tokenType tokens.TokenType, text []byte) []byte {
	regex, exists := tpm.compiledRegexMap[tokenType]
	if !exists {
		fmt.Printf("No regex found for TokenType %s\n", tokenType)
		return nil
	}
	return regex.Find(text)
}
