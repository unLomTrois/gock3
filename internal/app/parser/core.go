// core.go
package parser

import (
	"fmt"
	"strconv"

	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

// Expect verifies that the current token matches one of the expected types.
// If it does, it consumes the token and returns it.
// If not, it reports an error, attempts to recover, and returns nil.
func (p *Parser) Expect(expectedTypes ...tokens.TokenType) *tokens.Token {
	token := p.currentToken

	if token == nil {
		errMsg := fmt.Sprintf(errUnexpectedEOF, formatTokenTypes(expectedTypes))
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	// Check if current token matches any expected type
	for _, expectedType := range expectedTypes {
		if token.Type == expectedType {
			p.nextToken()
			return token
		}
	}

	// Token didn't match - report error and try to recover
	errMsg := fmt.Sprintf(errUnexpectedToken,
		token.Value,
		token.Type,
		formatTokenTypes(expectedTypes),
	)
	err := report.FromToken(token, severity.Error, errMsg)
	p.AddError(err)

	// Create a recovery point based on the expected types
	recoveryPoint := RecoveryPoint{
		TokenTypes: expectedTypes,
		Context:    "expected " + formatTokenTypes(expectedTypes),
	}

	// Attempt to recover
	if _, recovered := p.synchronize(recoveryPoint); recovered {
		return nil
	}

	// Recovery failed
	return nil
}

// unquoteExpect parses a quoted string, unquotes it, and returns the token.
func (p *Parser) unquoteExpect(expectedType tokens.TokenType) *tokens.Token {
	token := p.Expect(expectedType)
	if token == nil {
		return nil
	}

	unquotedValue, err := strconv.Unquote(token.Value)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedUnquoteString, token.Value)
		diag := report.FromToken(token, severity.Error, errMsg)
		p.AddError(diag)
		// Keep the original value if unquoting fails
		return token
	}

	token.Value = unquotedValue
	return token
}
