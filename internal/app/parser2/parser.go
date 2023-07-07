package parser2

import "ck3-parser/internal/app/lexer"

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

func (p *Parser) Parse() {
	p.lookahead = p.tokenstream.Next()
}

func (p *Parser) Expect(expectedtype lexer.TokenType) *lexer.Token {
	token := p.lookahead

	if token == nil {
		panic("[Parser] Unexpected end of input, expected: " + string(expectedtype))
	}
	if token.Type != expectedtype {
		panic("[Parser] Unexpected token: \"" + string(token.Value) + "\" with type of " + string(token.Type) + ", expected type: " + string(expectedtype))
	}

	return p.tokenstream.Next()
}
