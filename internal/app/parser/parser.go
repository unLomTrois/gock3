package parser

import (
	"ck3-parser/internal/app/lexer"
	"ck3-parser/internal/app/tokens"
	"fmt"
	"strconv"
)

type Parser struct {
	tokenstream *lexer.TokenStream
	lookahead   *tokens.Token
}

func New(tokenstream *lexer.TokenStream) *Parser {
	return &Parser{
		tokenstream: tokenstream,
		lookahead:   nil,
	}
}

func (p *Parser) Parse() []*Node {
	p.lookahead = p.tokenstream.Next()

	return p.List()
}

func (p *Parser) List(stop_lookahead ...tokens.TokenType) []*Node {
	nodes := make([]*Node, 0)

	for {
		if p.lookahead == nil {
			break
		}
		if len(stop_lookahead) > 0 && p.lookahead.Type == stop_lookahead[0] {
			break
		}

		new_node := p.Node()
		nodes = append(nodes, new_node)
	}

	return nodes
}

func (p *Parser) Node() *Node {

	switch p.lookahead.Type {
	case tokens.COMMENT:
		return p.CommentNode()
	default:
		return p.ExpressionNode()
	}
}

func (p *Parser) CommentNode() *Node {
	token := p.Expect(tokens.COMMENT)
	return &Node{
		Type:  Comment,
		Value: token.Value,
	}
}

func (p *Parser) ExpressionNode() *Node {
	key := p.Literal()

	var nodetype NodeType
	var operator *tokens.Token
	switch p.lookahead.Type {
	case tokens.EQUALS:
		operator = p.Expect(tokens.EQUALS)
		nodetype = Property
	case tokens.COMPARISON:
		operator = p.Expect(tokens.COMPARISON)
		nodetype = Comparison
	}

	switch p.lookahead.Type {
	case tokens.WORD, tokens.STRING, tokens.NUMBER, tokens.BOOL:
		value := p.Literal()
		node := &Node{
			Type:  nodetype,
			Key:   key,
			Value: value,
		}
		if nodetype == Comparison {
			node.Operator = operator.Value
		}
		return node
	case tokens.START:
		p.Expect(tokens.START)
		value := p.List(tokens.END)
		p.Expect(tokens.END)

		return &Node{
			Type:  Block,
			Key:   key,
			Value: value,
		}
	}

	return nil
}

func (p *Parser) Literal() *Literal {
	switch p.lookahead.Type {
	case tokens.WORD:
		return p.WordLiteral()
	case tokens.NUMBER:
		return p.NumberLiteral()
	case tokens.STRING:
		return p.StringLiteral()
	case tokens.BOOL:
		return p.BoolLiteral()
	default:
		panic(fmt.Sprintf("[Parser] Unexpected Literal: %q, with type of: %s",
			p.lookahead.Value, p.lookahead.Type))
	}
}

func (p *Parser) WordLiteral() *Literal {
	token := p.Expect(tokens.WORD)
	return &Literal{
		Type:  WordLiteral,
		Value: token.Value,
	}
}

func (p *Parser) NumberLiteral() *Literal {
	token := p.Expect(tokens.NUMBER)
	value, err := strconv.ParseFloat(token.Value, 32)
	if err != nil {
		panic(err)
	}

	return &Literal{
		Type:  NumberLiteral,
		Value: value,
	}
}

func (p *Parser) BoolLiteral() *Literal {
	token := p.Expect(tokens.BOOL)
	return &Literal{
		Type:  BoolLiteral,
		Value: token.Value,
	}
}

func (p *Parser) StringLiteral() *Literal {
	token := p.Expect(tokens.STRING)
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
