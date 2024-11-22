package parser

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

// Parser represents the parser with its current state and error manager.
type Parser struct {
	tokenstream *tokens.TokenStream
	lookahead   *tokens.Token
	loc         *tokens.Loc
	*report.ErrorManager
}

// New creates a new Parser instance.
func New(tokenstream *tokens.TokenStream) *Parser {
	return &Parser{
		tokenstream:  tokenstream,
		lookahead:    nil,
		loc:          nil,
		ErrorManager: report.NewErrorManager(),
	}
}

// Parse processes the token stream and returns the AST along with any diagnostic errors.
func Parse(token_stream *tokens.TokenStream) (*ast.FileBlock, []*report.DiagnosticItem) {
	p := New(token_stream)

	p.lookahead = p.tokenstream.Next()
	if p.lookahead != nil {
		p.loc = &p.lookahead.Loc
	}

	fileBlock := p.fileBlock()

	return fileBlock, p.Errors()
}

// fileBlock parses the entire file and constructs the AST's FileBlock.
func (p *Parser) fileBlock() *ast.FileBlock {
	if p.lookahead == nil {
		// Empty file
		return &ast.FileBlock{Values: []*ast.Field{}, Loc: tokens.Loc{}}
	}

	loc := p.loc
	fields := p.FieldList()
	return &ast.FileBlock{Values: fields, Loc: *loc}
}

// FieldList parses a list of fields until a stop token is encountered.
func (p *Parser) FieldList(stopLookahead ...tokens.TokenType) []*ast.Field {
	fields := make([]*ast.Field, 0)

	for p.lookahead != nil {
		// Check for stop tokens to end the field list
		if len(stopLookahead) > 0 && p.lookahead.Type == stopLookahead[0] {
			break
		}

		switch p.lookahead.Type {
		case tokens.COMMENT, tokens.NEXTLINE:
			p.Expect(tokens.COMMENT, tokens.NEXTLINE)
			continue
		case tokens.WORD, tokens.DATE, tokens.NUMBER:
			field := p.Field()
			if field != nil {
				fields = append(fields, field)
			}
		default:
			// Handle unexpected token
			errMsg := fmt.Sprintf(errFieldListUnexpectedToken, p.lookahead.Value, p.lookahead.Type)
			err := report.FromToken(p.lookahead, severity.Error, errMsg)
			p.AddError(err)

			if rec, recovered := p.synchronize(FieldListRecovery); !recovered {
				return fields // Stop parsing if recovery fails
			} else {
				log.Println(recovered, rec.Type, rec.Value)
			}
		}
	}

	return fields
}

// Field parses a single field and returns the corresponding AST node.
func (p *Parser) Field() *ast.Field {
	switch p.lookahead.Type {
	case tokens.WORD, tokens.DATE, tokens.NUMBER:
		return p.ExpressionNode()
	default:
		errMsg := fmt.Sprintf(errUnexpectedFieldToken, p.lookahead.Value, p.lookahead.Type)
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)

		if _, recovered := p.synchronize(FieldRecovery); !recovered {
			return nil // Stop parsing if recovery fails
		}
		return nil // Return nil after synchronization
	}
}

// ExpressionNode parses an expression node and returns the corresponding AST node.
func (p *Parser) ExpressionNode() *ast.Field {
	key := p.Key()
	if key == nil {
		return nil
	}

	operator := p.Operator()
	if operator == nil {
		return nil
	}

	value := p.Value()
	if value == nil {
		return nil
	}

	return &ast.Field{
		Key:      key,
		Operator: operator,
		Value:    value,
	}
}

// Key parses the key of a field and returns the corresponding token.
func (p *Parser) Key() *tokens.Token {
	if p.lookahead == nil {
		errMsg := "Expected a key, but reached end of input"
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	switch p.lookahead.Type {
	case tokens.WORD, tokens.DATE, tokens.NUMBER:
		return p.Expect(tokens.WORD, tokens.DATE, tokens.NUMBER)
	default:
		errMsg := fmt.Sprintf("Expected a key (WORD, DATE, or NUMBER), but found %q of type %q", p.lookahead.Value, p.lookahead.Type)
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)

		if _, recovered := p.synchronize(KeyRecovery); !recovered {
			return nil
		}
		return nil // Return nil after synchronization
	}
}

// Operator parses the operator of a field and returns the corresponding token.
func (p *Parser) Operator() *tokens.Token {
	if p.lookahead == nil {
		errMsg := errOperatorExpectedEOF
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	switch p.lookahead.Type {
	case tokens.QUESTION_EQUALS:
		return p.Expect(tokens.QUESTION_EQUALS)
	case tokens.EQUALS:
		return p.Expect(tokens.EQUALS)
	case tokens.COMPARISON:
		return p.Expect(tokens.COMPARISON)

	default:
		errMsg := fmt.Sprintf(errOperatorUnexpectedToken, p.lookahead.Value, p.lookahead.Type)
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)

		if _, recovered := p.synchronize(ValueRecovery); !recovered {
			return nil // Stop parsing if recovery fails
		}

		return nil // Return nil after synchronization, even if successful
	}
}

// Value parses the value of a field and returns the corresponding AST node.
func (p *Parser) Value() ast.BV {
	if p.lookahead == nil {
		errMsg := errValueExpectedEOF
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	switch p.lookahead.Type {
	case tokens.NEXTLINE:
		p.Expect(tokens.NEXTLINE)
		return p.EmptyValue()
	case tokens.WORD, tokens.NUMBER, tokens.QUOTED_STRING, tokens.BOOL:
		return p.Literal()
	case tokens.START:
		return p.Block()
	default:
		errMsg := fmt.Sprintf(errValueUnexpectedToken, p.lookahead.Value, p.lookahead.Type)
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize(ValueRecovery)
		return nil
	}
}

// Block parses a block and returns the corresponding AST node.
func (p *Parser) Block() ast.Block {
	p.Expect(tokens.START)
	loc := *p.loc

	if p.lookahead.Type == tokens.END {
		p.Expect(tokens.END)
		return &ast.FieldBlock{Values: []*ast.Field{}, Loc: loc}
	}

	var block ast.Block

	for p.lookahead.Type != tokens.END {
		switch p.lookahead.Type {
		case tokens.COMMENT, tokens.NEXTLINE:
			for p.lookahead.Type == tokens.COMMENT {
				p.Expect(tokens.COMMENT)
			}
			for p.lookahead.Type == tokens.NEXTLINE {
				p.Expect(tokens.NEXTLINE)
			}
			fallthrough
		case tokens.WORD, tokens.DATE:
			peek := p.tokenstream.Peek()
			if peek.Type != tokens.EQUALS && peek.Type != tokens.QUESTION_EQUALS {
				block = p.TokenBlock()
				break
			}

			block = p.FieldBlock(loc)
		case tokens.NUMBER, tokens.QUOTED_STRING:
			if p.tokenstream.Peek().Type == tokens.EQUALS {
				block = p.FieldBlock(loc)
				break
			}

			block = p.TokenBlock()
		default:
			errorMsg := fmt.Sprintf(errBlockUnexpectedToken, p.lookahead.Value, p.lookahead.Type)
			err := report.FromToken(p.lookahead, severity.Error, errorMsg)
			p.AddError(err)
			p.synchronize(BlockRecovery)
			continue
		}
	}

	// Expect closing brace '}'
	p.Expect(tokens.END)

	return block
}

// FieldBlock parses a block of fields and returns the corresponding AST node.
func (p *Parser) FieldBlock(loc tokens.Loc) *ast.FieldBlock {
	fields := p.FieldList(tokens.END)
	return &ast.FieldBlock{Values: fields, Loc: loc}
}

// TokenBlock parses a block of tokens and returns the corresponding AST node.
func (p *Parser) TokenBlock() *ast.TokenBlock {
	tokensList := p.TokenList(tokens.END)
	return &ast.TokenBlock{Values: tokensList}
}

// TokenList parses a list of tokens until a stop token is encountered.
func (p *Parser) TokenList(stopLookahead ...tokens.TokenType) []*tokens.Token {
	tokensList := make([]*tokens.Token, 0)

	for p.lookahead != nil {
		// Check for stop tokens to end the token list
		if len(stopLookahead) > 0 && p.lookahead.Type == stopLookahead[0] {
			break
		}

		switch p.lookahead.Type {
		case tokens.NEXTLINE:
			p.Expect(tokens.NEXTLINE)
			continue
		case tokens.NUMBER, tokens.QUOTED_STRING, tokens.WORD:
			token := p.Literal()
			if token != nil {
				tokensList = append(tokensList, token)
			}
		default:
			errMsg := fmt.Sprintf(errTokenListUnexpectedToken, p.lookahead.Value, p.lookahead.Type)
			err := report.FromToken(p.lookahead, severity.Error, errMsg)
			p.AddError(err)

			recoveryPoint := RecoveryPoint{
				TokenTypes: []tokens.TokenType{tokens.END, tokens.WORD, tokens.DATE},
				Context:    "TokenList",
			}

			if _, recovered := p.synchronize(recoveryPoint); !recovered {
				return tokensList // Stop parsing if recovery fails
			}
		}
	}

	return tokensList
}

// Literal parses a literal token and returns the corresponding token.
func (p *Parser) Literal() *tokens.Token {
	if p.lookahead == nil {
		err := report.FromLoc(*p.loc, severity.Error, errLiteralExpectedEOF)
		p.AddError(err)
		return nil
	}

	switch p.lookahead.Type {
	case tokens.WORD, tokens.NUMBER, tokens.BOOL:
		if token := p.Expect(p.lookahead.Type); token != nil {
			return token
		}
		// If Expect failed (shouldn't normally happen), try recovery

	case tokens.QUOTED_STRING:
		if token := p.unquoteExpect(tokens.QUOTED_STRING); token != nil {
			return token
		}
		// If unquoteExpect failed, try recovery

	default:
		errMsg := fmt.Sprintf(errLiteralUnexpectedToken, p.lookahead.Value, p.lookahead.Type)
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)
	}

	// Attempt recovery for any failure case
	if token, recovered := p.synchronize(LiteralRecovery); recovered {
		p.lookahead = token
		// Try parsing literal again from recovery point
		// But only try once to avoid potential infinite recursion
		if isLiteralType(token.Type) {
			return p.Literal()
		}
		// If recovered to a non-literal token, give up
		errMsg := fmt.Sprintf(errRecoveredNonLiteralToken, token.Value, token.Type)
		err := report.FromToken(token, severity.Error, errMsg)
		p.AddError(err)
	}

	return nil
}

// isLiteralType checks if a token type represents a literal value
func isLiteralType(tokenType tokens.TokenType) bool {
	switch tokenType {
	case tokens.WORD, tokens.NUMBER, tokens.BOOL, tokens.QUOTED_STRING:
		return true
	default:
		return false
	}
}

func (p *Parser) EmptyValue() ast.BV {
	return ast.EmptyValue{
		Loc: *p.loc,
	}
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

// Expect verifies that the current token matches one of the expected types.
// If it does, it consumes the token and returns it.
// If not, it reports an error, attempts to recover, and returns nil.
func (p *Parser) Expect(expectedTypes ...tokens.TokenType) *tokens.Token {
	token := p.lookahead

	if token == nil {
		errMsg := fmt.Sprintf(errUnexpectedEOF, formatTokenTypes(expectedTypes))
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	// Check if current token matches any expected type
	for _, expectedType := range expectedTypes {
		if token.Type == expectedType {
			p.loc = &token.Loc
			p.lookahead = p.tokenstream.Next()
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
	if nextToken, recovered := p.synchronize(recoveryPoint); recovered {
		p.lookahead = nextToken
		// Don't recursively call Expect - just return nil to indicate failure
		// but successful recovery
		return nil
	}

	// Recovery failed
	return nil
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
