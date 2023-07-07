package parser2

import (
	"ck3-parser/internal/app/lexer"
	"fmt"
)

type Parser struct {
	tokenstream *lexer.TokenStream
	lookahead   *lexer.Token
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

func (p *Parser) List(stop_lookahead ...lexer.TokenType) []*Node {
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
	// case lexer.SCRIPT:
	// 	return p.ScriptNode()
	case lexer.COMMENT:
		return p.CommentNode()
	default:
		return p.ExpressionNode()
	}
}

func (p *Parser) CommentNode() *Node {
	return &Node{
		Type:  Comment,
		Value: p.CommentLiteral(),
	}
}

func (p *Parser) ExpressionNode() *Node {
	key := p.Literal()

	var tokentype NodeType
	var operator *lexer.Token
	switch p.lookahead.Type {
	case lexer.EQUALS:
		operator = p.Expect(lexer.EQUALS)
		tokentype = Property
	case lexer.COMPARISON:
		operator = p.Expect(lexer.COMPARISON)
		tokentype = Comparison
	}

	switch p.lookahead.Type {
	case lexer.WORD, lexer.STRING, lexer.NUMBER, lexer.BOOL:
		value := p.Literal()
		node := &Node{
			Type:  tokentype,
			Key:   key,
			Value: value,
		}
		if tokentype == Comparison {
			node.Operator = operator.Value
		}
		return node
	case lexer.START:
		p.Expect(lexer.START)
		value := p.List(lexer.END)
		p.Expect(lexer.END)

		return &Node{
			Type:  Block,
			Key:   key,
			Value: value,
		}
	}

	return nil
}

func (p *Parser) Literal() interface{} {
	switch p.lookahead.Type {
	// case lexer.SCRIPT:
	// 	return p.ScriptLiteral()
	case lexer.WORD:
		return p.WordLiteral()
	case lexer.NUMBER:
		return p.NumberLiteral()
	case lexer.STRING:
		return p.StringLiteral()
	case lexer.BOOL:
		return p.BoolLiteral()
	case lexer.COMMENT:
		return p.CommentLiteral()
	default:
		panic(fmt.Sprintf("[Parser] Unexpected Literal: %q, with type of: %s",
			p.lookahead.Value, p.lookahead.Type))
	}
}

func (p *Parser) CommentLiteral() string {
	token := p.Expect(lexer.COMMENT)
	return token.Value
}

func (p *Parser) WordLiteral() string {
	token := p.Expect(lexer.WORD)
	return token.Value
}

func (p *Parser) NumberLiteral() string {
	token := p.Expect(lexer.NUMBER)
	return token.Value
}

func (p *Parser) BoolLiteral() string {
	token := p.Expect(lexer.BOOL)
	return token.Value
}

func (p *Parser) StringLiteral() string {
	token := p.Expect(lexer.STRING)
	return token.Value
}

// checks if the next token is the expected type and returns it
func (p *Parser) Expect(expectedtype lexer.TokenType) *lexer.Token {
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
