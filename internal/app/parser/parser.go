// internal/app/parser/parser.go
package parser

import (
	"strconv"

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
		case tokens.COMMENT:
			p.Expect(tokens.COMMENT)
			continue
		case tokens.WORD, tokens.DATE, tokens.NUMBER:
			field := p.Field()
			if field != nil {
				fields = append(fields, field)
			}
		default:
			// Handle unexpected token
			errMsg := "[FieldList] Unexpected token '" + p.lookahead.Value + "' of type '" + string(p.lookahead.Type) + "'"
			err := report.FromToken(p.lookahead, severity.Error, errMsg)
			p.AddError(err)
			p.synchronize(tokens.END, tokens.WORD, tokens.DATE)
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
		errMsg := "Unexpected token '" + p.lookahead.Value + "' of type '" + string(p.lookahead.Type) + "' when expecting a field"
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)
		p.lookahead = p.tokenstream.Next() // Advance to recover
		return nil
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
		errMsg := "Expected a key (WORD or DATE), but found '" + p.lookahead.Value + "' of type '" + string(p.lookahead.Type) + "'"
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize(tokens.EQUALS, tokens.COMPARISON, tokens.END, tokens.WORD)
		return nil
	}
}

// Operator parses the operator of a field and returns the corresponding token.
func (p *Parser) Operator() *tokens.Token {
	if p.lookahead == nil {
		errMsg := "Expected an operator '=', '==', or comparison, but reached end of input"
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	// QUESTION_EQUALS

	switch p.lookahead.Type {
	case tokens.QUESTION_EQUALS:
		return p.Expect(tokens.QUESTION_EQUALS)
	case tokens.EQUALS:
		return p.Expect(tokens.EQUALS)
	case tokens.COMPARISON:
		return p.Expect(tokens.COMPARISON)

	default:
		errMsg := "Expected operator '=', '==', or comparison, but found '" + p.lookahead.Value + "' of type '" + string(p.lookahead.Type) + "'"
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize(tokens.WORD, tokens.NUMBER, tokens.QUOTED_STRING, tokens.BOOL, tokens.START)
		return nil
	}
}

// Value parses the value of a field and returns the corresponding AST node.
func (p *Parser) Value() ast.BV {
	if p.lookahead == nil {
		errMsg := "Expected a value, but reached end of input"
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	switch p.lookahead.Type {
	case tokens.WORD, tokens.NUMBER, tokens.QUOTED_STRING, tokens.BOOL:
		return p.Literal()
	case tokens.START:
		return p.Block()
	default:
		errMsg := "Unexpected token '" + p.lookahead.Value + "' of type '" + string(p.lookahead.Type) + "' in value"
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize(tokens.END, tokens.WORD, tokens.DATE)
		return nil
	}
}

// Block parses a block and returns the corresponding AST node.
func (p *Parser) Block() ast.Block {
	if p.lookahead.Type != tokens.START {
		errMsg := "Expected '{', but found '" + p.lookahead.Value + "' of type '" + string(p.lookahead.Type) + "'"
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize(tokens.START, tokens.END, tokens.WORD, tokens.DATE)
		return nil
	}

	p.Expect(tokens.START)
	loc := *p.loc

	if p.lookahead.Type == tokens.END {
		p.Expect(tokens.END)
		return &ast.FieldBlock{Values: []*ast.Field{}, Loc: loc}
	}

	var block ast.Block

	switch p.lookahead.Type {
	case tokens.COMMENT:
		p.Expect(tokens.COMMENT)
		fallthrough
	case tokens.WORD, tokens.DATE:
		peek := p.tokenstream.Peek()
		if peek.Type != tokens.EQUALS && peek.Type != tokens.QUESTION_EQUALS {
			block = p.TokenBlock()
			break
		}

		block = p.FieldBlock(loc)
		// case tokens.QUOTED_STRING:
	case tokens.NUMBER, tokens.QUOTED_STRING:
		if p.tokenstream.Peek().Type == tokens.EQUALS {
			block = p.FieldBlock(loc)
			break
		}

		block = p.TokenBlock()
	default:
		errMsg := "[Block] Unexpected token '" + p.lookahead.Value + "' of type '" + string(p.lookahead.Type) + "' in block"
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize(tokens.END, tokens.WORD, tokens.DATE)
		return nil
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
		case tokens.NUMBER, tokens.QUOTED_STRING, tokens.WORD:
			token := p.Literal()
			if token != nil {
				tokensList = append(tokensList, token)
			}
		default:
			// Handle unexpected token
			errMsg := "[TokenList] Unexpected token '" + p.lookahead.Value + "' of type '" + string(p.lookahead.Type) + "' in token list"
			err := report.FromToken(p.lookahead, severity.Error, errMsg)
			p.AddError(err)
			p.synchronize(tokens.END, tokens.WORD, tokens.DATE)
		}
	}

	return tokensList
}

// Literal parses a literal token and returns the corresponding token.
func (p *Parser) Literal() *tokens.Token {
	switch p.lookahead.Type {
	case tokens.WORD, tokens.NUMBER, tokens.BOOL:
		return p.Expect(p.lookahead.Type)
	case tokens.QUOTED_STRING:
		return p.unquoteExpect(tokens.QUOTED_STRING)
	default:
		errMsg := "Unexpected token '" + p.lookahead.Value + "' of type '" + string(p.lookahead.Type) + "' as literal"
		err := report.FromToken(p.lookahead, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize()
		return nil
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
		errMsg := "Failed to unquote string '" + token.Value + "'"
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
		errMsg := "Unexpected end of input, expected one of: " + formatTokenTypes(expectedTypes)
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	for _, expectedType := range expectedTypes {
		if token.Type == expectedType {
			p.loc = &token.Loc
			p.lookahead = p.tokenstream.Next()
			return token
		}
	}

	// Report unexpected token
	errMsg := "Unexpected token '" + token.Value + "' of type '" + string(token.Type) + "', expected one of: " + formatTokenTypes(expectedTypes)
	err := report.FromToken(token, severity.Error, errMsg)
	p.AddError(err)

	// Attempt to recover by synchronizing
	p.synchronize(expectedTypes...)

	return nil
}

// synchronize advances the parser's lookahead until it finds a synchronization token or one of the expected types.
func (p *Parser) synchronize(expectedTypes ...tokens.TokenType) {
	for p.lookahead != nil {
		// If the current token matches any of the expected types, stop synchronizing
		for _, expectedType := range expectedTypes {
			if p.lookahead.Type == expectedType {
				p.Expect(expectedType)
			}
		}

		// Otherwise, consume the token and continue synchronizing
		p.lookahead = p.tokenstream.Next()
	}
}

// formatTokenTypes formats a slice of TokenType into a readable string.
func formatTokenTypes(types []tokens.TokenType) string {
	formatted := ""
	for i, t := range types {
		formatted += "'" + string(t) + "'"
		if i < len(types)-1 {
			formatted += ", "
		}
	}
	return formatted
}
