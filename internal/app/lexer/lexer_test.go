package lexer

import (
	"ck3-parser/internal/app/tokens"
	"reflect"
	"testing"
)

const (
	elementary = `
  namespace = cooking

  entity = {
    scope:character = character.123
  }
`
)

func TestLexer_GetNextToken(t *testing.T) {
	lexer := NewLexer([]byte(elementary))

	tests := []struct {
		name     string
		want     *tokens.Token
		wantErr  bool
		skipNext bool
	}{
		{
			name: "Namespace is WORD",
			want: &tokens.Token{Type: tokens.WORD, Value: "namespace"},
		},
		{
			name: "= is EQUAL",
			want: &tokens.Token{Type: tokens.EQUALS, Value: "="},
		},
		{
			name: "cooking is WORD",
			want: &tokens.Token{Type: tokens.WORD, Value: "cooking"},
		},
		{
			name:     "entity is WORD",
			want:     &tokens.Token{Type: tokens.WORD, Value: "entity"},
			skipNext: true,
		},
		{
			name: "{ is START",
			want: &tokens.Token{Type: tokens.START, Value: "{"},
		},
		{
			name:     "scope:character is WORD",
			want:     &tokens.Token{Type: tokens.WORD, Value: "scope:character"},
			skipNext: true,
		},
		{
			name:     "character.123 is WORD",
			want:     &tokens.Token{Type: tokens.WORD, Value: "character.123"},
			skipNext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lexer.getNextToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Lexer.getNextToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.getNextToken() = %v, want %v", got, tt.want)
			}
			if tt.skipNext {
				_, _ = lexer.getNextToken() // Ignore errors for skipped tokens
			}
		})
	}
}

func TestLexer_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []*tokens.Token
		wantErr bool
	}{
		{
			name:  "Elementary is tokenized correctly",
			input: elementary,
			want: []*tokens.Token{
				{Type: tokens.WORD, Value: "namespace"},
				{Type: tokens.EQUALS, Value: "="},
				{Type: tokens.WORD, Value: "cooking"},
				{Type: tokens.WORD, Value: "entity"},
				{Type: tokens.EQUALS, Value: "="},
				{Type: tokens.START, Value: "{"},
				{Type: tokens.WORD, Value: "scope:character"},
				{Type: tokens.EQUALS, Value: "="},
				{Type: tokens.WORD, Value: "character.123"},
				{Type: tokens.END, Value: "}"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input))

			got, err := lexer.Scan()
			if (err != nil) != tt.wantErr {
				t.Errorf("Lexer.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Stream, tt.want) {
				t.Errorf("Lexer.Scan() =\ngot:  %v\nwant: %v", got.Stream, tt.want)
			}
		})
	}
}
