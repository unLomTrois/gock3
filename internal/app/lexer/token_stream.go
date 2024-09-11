package lexer

import "ck3-parser/internal/app/tokens"

type TokenStream struct {
	Stream []*tokens.Token
	Cursor int
}

func NewTokenStream() *TokenStream {
	return &TokenStream{
		Stream: []*tokens.Token{},
		Cursor: 0,
	}
}

func (ts *TokenStream) Push(token *tokens.Token) *TokenStream {
	ts.Stream = append(ts.Stream, token)

	return ts
}

func (ts *TokenStream) Next() *tokens.Token {
	if ts.Cursor >= len(ts.Stream) {
		return nil
	}
	token := ts.Stream[ts.Cursor]
	ts.Cursor++

	return token
}
