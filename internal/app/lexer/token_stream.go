package lexer

type TokenStream struct {
	Stream []*Token
	Cursor int
}

func NewTokenStream() *TokenStream {
	return &TokenStream{
		Stream: []*Token{},
		Cursor: 0,
	}
}

func (ts *TokenStream) Push(token *Token) *TokenStream {
	ts.Stream = append(ts.Stream, token)

	return ts
}

func (ts *TokenStream) Next() *Token {
	if ts.Cursor >= len(ts.Stream) {
		return nil
	}
	token := ts.Stream[ts.Cursor]
	ts.Cursor++

	return token
}
