package tokens

type TokenStream struct {
	tokens   []*Token
	position int
}

func NewTokenStream() *TokenStream {
	return &TokenStream{
		tokens:   []*Token{},
		position: 0,
	}
}

func (ts *TokenStream) Push(token *Token) *TokenStream {
	ts.tokens = append(ts.tokens, token)

	return ts
}

func (ts *TokenStream) Next() *Token {
	if ts.position >= len(ts.tokens) {
		return nil
	}
	token := ts.tokens[ts.position]
	ts.position++

	return token
}

func (ts *TokenStream) Peek() *Token {
	if ts.position < len(ts.tokens) {
		return ts.tokens[ts.position]
	}
	return nil
}
