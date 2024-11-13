package tokens

type TokenStream struct {
	Tokens   []*Token
	Position int
}

func NewTokenStream() *TokenStream {
	return &TokenStream{
		Tokens:   []*Token{},
		Position: 0,
	}
}

func (ts *TokenStream) Push(token *Token) *TokenStream {
	ts.Tokens = append(ts.Tokens, token)

	return ts
}

func (ts *TokenStream) Next() *Token {
	if ts.Position >= len(ts.Tokens) {
		return nil
	}
	token := ts.Tokens[ts.Position]
	ts.Position++

	return token
}

func (ts *TokenStream) Peek() *Token {
	if ts.Position < len(ts.Tokens) {
		return ts.Tokens[ts.Position]
	}
	return nil
}
