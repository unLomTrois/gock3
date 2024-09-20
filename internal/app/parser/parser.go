package parser

import (
	"fmt"

	"github.com/unLomTrois/lexiCK3/internal/app/tokens"
)

type Parser struct {
	tokenstream *tokens.TokenStream
	lookahead   *tokens.Token
}

func New(tokenstream *tokens.TokenStream) *Parser {
	return &Parser{
		tokenstream: tokenstream,
		lookahead:   nil,
	}
}

func Parse(token_stream *tokens.TokenStream) *FileBlock {
	p := New(token_stream)

	p.lookahead = p.tokenstream.Next()

	return p.fileBlock()
}

func (p *Parser) fileBlock() *FileBlock {
	return &FileBlock{Values: p.FieldList()}
}

func (p *Parser) FieldList(stop_lookahead ...tokens.TokenType) []*Field {
	fields := make([]*Field, 0)

	for p.lookahead != nil {
		if len(stop_lookahead) > 0 && p.lookahead.Type == stop_lookahead[0] {
			break
		}

		switch p.lookahead.Type {
		case tokens.COMMENT:
			p.Expect(tokens.COMMENT)
			continue
		case tokens.WORD:
			field := p.Field()
			fields = append(fields, field)
		case tokens.DATE:
			field := p.Field()
			fields = append(fields, field)
		default:
			// If the current symbol is not in FIRST(Statement), then it is an ε-production
			panic(fmt.Sprintf("[Parser] Unexpected Statement: %q, with type of: %s",
				p.lookahead.Value, p.lookahead.Type))
		}
	}

	return fields
}

func (p *Parser) Field() *Field {
	switch p.lookahead.Type {
	case tokens.WORD, tokens.DATE:
		return p.ExpressionNode()
	default:
		panic(fmt.Sprintf("[Parser] Unexpected Node: %q, with type of: %s",
			p.lookahead.Value, p.lookahead.Type))
	}
}

func (p *Parser) ExpressionNode() *Field {
	key := p.Key()

	operator, err := p.Operator()
	if err != nil {
		panic(err)
	}

	value, err := p.Value()
	if err != nil {
		panic(err)
	}

	return &Field{
		Key:      key,
		Operator: operator,
		Value:    value,
	}
}

func (p *Parser) Key() *tokens.Token {
	switch p.lookahead.Type {
	case tokens.WORD:
		return p.Expect(tokens.WORD)
	case tokens.DATE:
		return p.Expect(tokens.DATE)
	default:
		panic(fmt.Sprintf("[Parser] Unexpected Key: %q, with type of: %s",
			p.lookahead.Value, p.lookahead.Type))
	}
}

func (p *Parser) Operator() (*tokens.Token, error) {
	switch p.lookahead.Type {
	case tokens.EQUALS, tokens.COMPARISON:
		return p.Expect(p.lookahead.Type), nil
	default:
		return nil, fmt.Errorf("expected operator '=', '==', or comparison, got %s", p.lookahead.Type)
	}
}

func (p *Parser) Value() (BV, error) {
	switch p.lookahead.Type {
	case tokens.WORD, tokens.NUMBER, tokens.QUOTED_STRING, tokens.BOOL:
		return p.Literal(), nil
	case tokens.START:
		return p.Block()
	default:
		return nil, fmt.Errorf("unexpected token %s in Value", p.lookahead.Type)
	}
}

func (p *Parser) Block() (Block, error) {
	p.Expect(tokens.START)

	switch p.lookahead.Type {
	case tokens.COMMENT:
		p.Expect(tokens.COMMENT)
		fallthrough
	case tokens.WORD:
		return p.FieldBlock(), nil
	case tokens.NUMBER, tokens.QUOTED_STRING:
		return p.TokenBlock(), nil
	default:
		return nil, fmt.Errorf("unexpected token %s in Block", p.lookahead.Type)
	}
}

func (p *Parser) FieldBlock() *FieldBlock {
	nodes := p.FieldList(tokens.END)
	p.Expect(tokens.END)
	return &FieldBlock{Values: nodes}
}

func (p *Parser) TokenBlock() *TokenBlock {
	nodes := p.TokenList(tokens.END)
	p.Expect(tokens.END)
	return &TokenBlock{Values: nodes}
}

func (p *Parser) TokenList(stop_lookahead ...tokens.TokenType) []*tokens.Token {
	nodes := make([]*tokens.Token, 0)

	for p.lookahead != nil {
		if len(stop_lookahead) > 0 && p.lookahead.Type == stop_lookahead[0] {
			break
		}

		switch p.lookahead.Type {
		case tokens.NUMBER, tokens.QUOTED_STRING:
			node := p.Literal()
			nodes = append(nodes, node)
		default:
			// If the current symbol is not in FIRST(Statement), then it is an ε-production
			panic(fmt.Sprintf("[Parser] Unexpected Statement: %q, with type of: %s",
				p.lookahead.Value, p.lookahead.Type))
		}
	}

	return nodes
}

func (p *Parser) Literal() *tokens.Token {
	switch p.lookahead.Type {
	case tokens.WORD:
		return p.Expect(tokens.WORD)
	case tokens.NUMBER:
		return p.Expect(tokens.NUMBER)
	case tokens.QUOTED_STRING:
		return p.Expect(tokens.QUOTED_STRING)
	case tokens.BOOL:
		return p.Expect(tokens.BOOL)
	default:
		panic(fmt.Sprintf("[Parser] Unexpected Literal: %q, with type of: %s",
			p.lookahead.Value, p.lookahead.Type))
	}
}

// checks if the next token is the expected type and returns it
func (p *Parser) Expect(expectedtype tokens.TokenType) *tokens.Token {
	token := p.lookahead

	if token == nil {
		panic("[Parser] Unexpected end of input, expected: " + string(expectedtype))
	}
	if token.Type != expectedtype {
		fmt.Println(p.tokenstream.Cursor)
		panic("[Parser] Unexpected token: \"" + string(token.Value) + "\" with type of " + string(token.Type) + "\nexpected type: " + string(expectedtype))
	}

	p.lookahead = p.tokenstream.Next()

	return token
}
