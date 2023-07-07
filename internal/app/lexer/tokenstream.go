package lexer

type TokenStream struct {
	Stream []*Token
	cursor int
}

func NewTokenStream() *TokenStream {
	return &TokenStream{
		Stream: []*Token{},
		cursor: 0,
	}
}

func (ts *TokenStream) Push(token *Token) *TokenStream {
	ts.Stream = append(ts.Stream, token)

	return ts
}

func (ts *TokenStream) Next() *Token {
	if ts.cursor >= len(ts.Stream) {
		return nil
	}
	token := ts.Stream[ts.cursor]
	ts.cursor++

	return token
}
