// block.go
package parser

import (
	"fmt"

	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

// Block parses a block and returns the corresponding AST node.
func (p *Parser) Block() ast.Block {
	p.Expect(tokens.START)
	loc := *p.loc

	if p.currentToken.Type == tokens.END {
		p.Expect(tokens.END)
		return &ast.FieldBlock{Values: []*ast.Field{}, Loc: loc}
	}

	var block ast.Block

	for p.currentToken != nil && p.currentToken.Type != tokens.END {
		switch p.currentToken.Type {
		case tokens.NEXTLINE:
			p.skipTokens(tokens.NEXTLINE)
			continue
		case tokens.WORD, tokens.DATE, tokens.QUOTED_STRING, tokens.NUMBER:
			if p.isNextField() {
				block = p.FieldBlock(loc)
			} else {
				block = p.TokenBlock()
			}
		default:
			errorMsg := fmt.Sprintf(errBlockUnexpectedToken, p.currentToken.Value, p.currentToken.Type)
			err := report.FromToken(p.currentToken, severity.Error, errorMsg)
			p.AddError(err)
			p.synchronize(BlockRecovery)
			continue
		}
		break // Exit the loop after processing the block
	}

	// Expect closing brace '}'
	p.Expect(tokens.END)

	return block
}

func (p *Parser) skipTokens(types ...tokens.TokenType) {
	for p.currentToken != nil {
		match := false
		for _, t := range types {
			match = p.currentToken.Type == t
			if match {
				p.Expect(t) // Consume the token
				break
			}
		}
		if !match {
			break // Exit if currentToken doesn't match any type
		}
	}
}

// isNextField determines if the next construct is likely a field.
func (p *Parser) isNextField() bool {
	return isKeyToken(p.currentToken.Type) && isOperatorToken(p.lookahead.Type)
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

	for p.currentToken != nil {
		// Check for stop tokens to end the token list
		if len(stopLookahead) > 0 && p.currentToken.Type == stopLookahead[0] {
			break
		}

		switch p.currentToken.Type {
		case tokens.NEXTLINE:
			p.Expect(tokens.NEXTLINE)
			continue
		case tokens.NUMBER, tokens.QUOTED_STRING, tokens.WORD:
			token := p.Literal()
			if token != nil {
				tokensList = append(tokensList, token)
			}
		default:
			errMsg := fmt.Sprintf(errTokenListUnexpectedToken, p.currentToken.Value, p.currentToken.Type)
			err := report.FromToken(p.currentToken, severity.Error, errMsg)
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
