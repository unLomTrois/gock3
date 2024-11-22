package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

// Recovery mode constants
const (
	maxRecoveryAttempts = 10
	maxTokensToSkip     = 50
)

// RecoveryPoint represents a stable point in the grammar where parsing can resume
type RecoveryPoint struct {
	TokenTypes []tokens.TokenType
	Context    string
}

// Common recovery points in the grammar - exported so they can be used in parser.go
var (
	// For field-level recovery - look for start of new field
	FieldRecovery = RecoveryPoint{
		TokenTypes: []tokens.TokenType{tokens.WORD, tokens.DATE, tokens.NUMBER, tokens.END},
		Context:    "field",
	}

	// For block-level recovery - look for block end or new statement
	BlockRecovery = RecoveryPoint{
		TokenTypes: []tokens.TokenType{tokens.END, tokens.WORD, tokens.DATE},
		Context:    "block",
	}

	// For expression recovery - look for operators or statement end
	ExpressionRecovery = RecoveryPoint{
		TokenTypes: []tokens.TokenType{tokens.EQUALS, tokens.COMPARISON, tokens.END},
		Context:    "expression",
	}

	KeyRecovery = RecoveryPoint{
		TokenTypes: []tokens.TokenType{tokens.EQUALS, tokens.COMPARISON, tokens.END, tokens.WORD},
		Context:    "key",
	}

	//
	ValueRecovery = RecoveryPoint{
		TokenTypes: []tokens.TokenType{tokens.WORD, tokens.NUMBER, tokens.QUOTED_STRING, tokens.BOOL, tokens.START},
		Context:    "value",
	}

	FieldListRecovery = RecoveryPoint{
		TokenTypes: []tokens.TokenType{
			tokens.END, tokens.WORD, tokens.DATE,
			// tokens.WORD, tokens.DATE, tokens.NUMBER, // Valid field starters
			// tokens.NEXTLINE, tokens.COMMENT, // Skippable tokens
			// tokens.END, // Block end
		},
		Context: "field list",
	}

	LiteralRecovery = RecoveryPoint{
		TokenTypes: []tokens.TokenType{
			tokens.WORD,
			tokens.NUMBER,
			tokens.BOOL,
			tokens.QUOTED_STRING,
			tokens.NEXTLINE, // Allow recovery at line boundaries
			tokens.END,      // Allow recovery at block ends
		},
		Context: "literal value",
	}
)

// synchronize attempts to recover from parsing errors by finding a stable point
func (p *Parser) synchronize(point RecoveryPoint) (*tokens.Token, bool) {
	attempts := 0
	skippedTokens := 0
	startLoc := *p.loc

	// Keep track of skipped tokens for error reporting
	var skipped []*tokens.Token

	for p.lookahead != nil && attempts < maxRecoveryAttempts && skippedTokens < maxTokensToSkip {
		attempts++

		log.Println(attempts)

		// Check if current token is a recovery point
		for _, expectedType := range point.TokenTypes {
			if p.lookahead.Type == expectedType {
				// Found recovery point - report skipped section
				if len(skipped) > 0 {
					p.reportSkippedSection(startLoc, skipped, point.Context)
				}
				return p.lookahead, true
			}
		}

		// Skip current token
		skipped = append(skipped, p.lookahead)
		p.lookahead = p.tokenstream.Next()
		skippedTokens++
	}

	// Failed to recover
	p.reportRecoveryFailure(startLoc, point.Context)
	return nil, false
}

// ? is it even needed?
// reportSkippedSection reports the tokens that were skipped during recovery
func (p *Parser) reportSkippedSection(startLoc tokens.Loc, skipped []*tokens.Token, context string) {
	var skippedValues []string
	for _, t := range skipped {
		skippedValues = append(skippedValues, fmt.Sprintf("%s (%s)", t.Value, t.Type))
	}

	errMsg := fmt.Sprintf(
		"Skipped invalid syntax in %q: %q",
		context,
		strings.Join(skippedValues, ", "),
	)

	err := report.FromLoc(startLoc, severity.Warning, errMsg)
	p.AddError(err)
}

// reportRecoveryFailure reports when the parser couldn't recover
func (p *Parser) reportRecoveryFailure(startLoc tokens.Loc, context string) {
	errMsg := fmt.Sprintf(
		"Failed to recover while parsing %s - too many invalid tokens",
		context,
	)
	err := report.FromLoc(startLoc, severity.Error, errMsg)
	p.AddError(err)
}

// Optional: Helper functions specific to error recovery
func formatSkippedTokens(tokens []*tokens.Token) string {
	var parts []string
	for _, t := range tokens {
		parts = append(parts, fmt.Sprintf("%q(%s)", t.Value, t.Type))
	}
	return strings.Join(parts, ", ")
}
