package lexer

import (
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

	scriptedTrigger = `
scripted_trigger cooking_trigger = {
  condition1 = yes
  condition2 = no
}
`
)

func TestLexer_GetNextToken(t *testing.T) {
	lexer := NewLexer([]byte(elementary))

	tests := []struct {
		name     string
		want     *Token
		wantErr  bool
		skipNext bool
	}{
		{
			name: "Namespace is WORD",
			want: &Token{Type: WORD, Value: "namespace"},
		},
		{
			name: "= is EQUAL",
			want: &Token{Type: EQUALS, Value: "="},
		},
		{
			name: "cooking is WORD",
			want: &Token{Type: WORD, Value: "cooking"},
		},
		{
			name:     "entity is WORD",
			want:     &Token{Type: WORD, Value: "entity"},
			skipNext: true,
		},
		{
			name: "{ is START",
			want: &Token{Type: START, Value: "{"},
		},
		{
			name:     "scope:character is WORD",
			want:     &Token{Type: WORD, Value: "scope:character"},
			skipNext: true,
		},
		{
			name:     "character.123 is WORD",
			want:     &Token{Type: WORD, Value: "character.123"},
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
		want    []*Token
		wantErr bool
	}{
		{
			name:  "Elementary is tokenized correctly",
			input: elementary,
			want: []*Token{
				{Type: WORD, Value: "namespace"},
				{Type: EQUALS, Value: "="},
				{Type: WORD, Value: "cooking"},
				{Type: WORD, Value: "entity"},
				{Type: EQUALS, Value: "="},
				{Type: START, Value: "{"},
				{Type: WORD, Value: "scope:character"},
				{Type: EQUALS, Value: "="},
				{Type: WORD, Value: "character.123"},
				{Type: END, Value: "}"},
			},
		},
		{
			name:  "Scripted trigger is tokenized correctly",
			input: scriptedTrigger,
			want: []*Token{
				{Type: SCRIPT, Value: "scripted_trigger"},
				{Type: WORD, Value: "cooking_trigger"},
				{Type: EQUALS, Value: "="},
				{Type: START, Value: "{"},
				{Type: WORD, Value: "condition1"},
				{Type: EQUALS, Value: "="},
				{Type: BOOL, Value: "yes"},
				{Type: WORD, Value: "condition2"},
				{Type: EQUALS, Value: "="},
				{Type: BOOL, Value: "no"},
				{Type: END, Value: "}"},
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
