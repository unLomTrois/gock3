// helpers.go
package parser

import (
	"fmt"
	"strings"

	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
)

// Helper functions for token type checks.
func isKeyToken(tokenType tokens.TokenType) bool {
	return tokenType == tokens.WORD || tokenType == tokens.DATE || tokenType == tokens.NUMBER
}

func isOperatorToken(tokenType tokens.TokenType) bool {
	return isEqualOperatorToken(tokenType) || tokenType == tokens.COMPARISON
}

func isEqualOperatorToken(tokenType tokens.TokenType) bool {
	return tokenType == tokens.EQUALS || tokenType == tokens.QUESTION_EQUALS
}

// isLiteralType checks if a token type represents a literal value.
func isLiteralType(tokenType tokens.TokenType) bool {
	switch tokenType {
	case tokens.WORD, tokens.NUMBER, tokens.BOOL, tokens.QUOTED_STRING:
		return true
	default:
		return false
	}
}

// formatTokenTypes formats a slice of TokenType into a readable string.
func formatTokenTypes(types []tokens.TokenType) string {
	if len(types) == 0 {
		return "no token types specified"
	}

	if len(types) == 1 {
		return fmt.Sprintf("%q", types[0])
	}

	parts := make([]string, len(types))
	for i, t := range types {
		parts[i] = fmt.Sprintf("%q", t)
	}

	// For two types use "x or y"
	if len(types) == 2 {
		return fmt.Sprintf("%s or %s", parts[0], parts[1])
	}

	// For more than two types use "x, y, or z"
	return fmt.Sprintf("%s, or %s",
		strings.Join(parts[:len(parts)-1], ", "),
		parts[len(parts)-1],
	)
}
