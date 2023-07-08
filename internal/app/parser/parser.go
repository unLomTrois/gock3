package parser

import (
	"ck3-parser/internal/app/lexer"
	"fmt"
	"strconv"
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
	token := p.Expect(lexer.COMMENT)
	return &Node{
		Type:  Comment,
		Value: token.Value,
	}
}

func (p *Parser) ExpressionNode() *Node {
	key := p.Literal()

	var nodetype NodeType
	var operator *lexer.Token
	switch p.lookahead.Type {
	case lexer.EQUALS:
		operator = p.Expect(lexer.EQUALS)
		nodetype = Property
	case lexer.COMPARISON:
		operator = p.Expect(lexer.COMPARISON)
		nodetype = Comparison
	}

	switch p.lookahead.Type {
	case lexer.WORD, lexer.STRING, lexer.NUMBER, lexer.BOOL:
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

func (p *Parser) Literal() *Literal {
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
	default:
		panic(fmt.Sprintf("[Parser] Unexpected Literal: %q, with type of: %s",
			p.lookahead.Value, p.lookahead.Type))
	}
}

func (p *Parser) WordLiteral() *Literal {
	token := p.Expect(lexer.WORD)
	return &Literal{
		Type:  WordLiteral,
		Value: token.Value,
	}
}

func (p *Parser) NumberLiteral() *Literal {
	token := p.Expect(lexer.NUMBER)
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
	token := p.Expect(lexer.BOOL)
	return &Literal{
		Type:  BoolLiteral,
		Value: token.Value,
	}
}

func (p *Parser) StringLiteral() *Literal {
	token := p.Expect(lexer.STRING)
	return &Literal{
		Type:  StringLiteral,
		Value: token.Value,
	}
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
