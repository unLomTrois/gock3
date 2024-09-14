package parser

import (
	"ck3-parser/internal/app/tokens"
	"fmt"
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

func Parse(token_stream *tokens.TokenStream) []*Node {
	p := New(token_stream)

	p.lookahead = p.tokenstream.Next()

	return p.Statements()
}

func (p *Parser) Statements(stop_lookahead ...tokens.TokenType) []*Node {
	nodes := make([]*Node, 0)

	for p.lookahead != nil {
		if len(stop_lookahead) > 0 && p.lookahead.Type == stop_lookahead[0] {
			break
		}

		switch p.lookahead.Type {
		case tokens.COMMENT, tokens.WORD:
			node := p.Node()
			nodes = append(nodes, node)
		default:
			// If the current symbol is not in FIRST(Statement), then it is an ε-production
		}
	}

	return nodes
}

func (p *Parser) Node() *Node {
	switch p.lookahead.Type {
	case tokens.COMMENT:
		return p.CommentNode()
	case tokens.WORD:
		return p.ExpressionNode()
	default:
		panic(fmt.Sprintf("[Parser] Unexpected Node: %q, with type of: %s",
			p.lookahead.Value, p.lookahead.Type))
	}
}

func (p *Parser) CommentNode() *Node {
	token := p.Expect(tokens.COMMENT)
	return &Node{
		Value: token,
	}
}

func (p *Parser) ExpressionNode() *Node {
	key := p.Key()

	operator, err := p.Operator()
	if err != nil {
		panic(err)
	}

	value, err := p.Value()
	if err != nil {
		panic(err)
	}

	return &Node{
		Key:      key,
		Operator: operator,
		Value:    value,
	}
}

func (p *Parser) Key() *tokens.Token {
	return p.Expect(tokens.WORD)
}

func (p *Parser) Operator() (*tokens.Token, error) {
	switch p.lookahead.Type {
	case tokens.EQUALS, tokens.COMPARISON:
		return p.Expect(p.lookahead.Type), nil
	default:
		return nil, fmt.Errorf("expected operator '=', '==', or comparison, got %s", p.lookahead.Type)
	}
}

func (p *Parser) Value() (interface{}, error) {
	switch p.lookahead.Type {
	case tokens.WORD, tokens.NUMBER, tokens.QUOTED_STRING, tokens.BOOL:
		return p.Literal(), nil
	case tokens.START:
		return p.Block()
	default:
		return nil, fmt.Errorf("unexpected token %s in Value", p.lookahead.Type)
	}
}

func (p *Parser) Block() ([]*Node, error) {
	p.Expect(tokens.START)
	nodes := p.Statements(tokens.END)
	p.Expect(tokens.END)
	return nodes, nil
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

func (p *Parser) WordLiteral() *tokens.Token {
	return p.Expect(tokens.WORD)
}

func (p *Parser) StringLiteral() *Literal {
	token := p.Expect(tokens.QUOTED_STRING)
	return &Literal{
		Type:  StringLiteral,
		Value: token.Value,
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
