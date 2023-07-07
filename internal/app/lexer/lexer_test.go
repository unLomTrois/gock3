package lexer

import (
	"reflect"
	"testing"
)

var elementary = `
namespace = cooking

entity = {
  scope:character = character.123
}
`

var scripted_trigger = `
scripted_trigger cooking_trigger = {
  condition1 = yes
  condition2 = no
}
`

func TestLexer_GetNextToken(t *testing.T) {
	rawtext := []byte(elementary)

	lexer := New(rawtext)

	tests := []struct {
		name       string
		want       *Token
		wantErr    bool
		posteffect func()
		skipnext   bool
	}{
		{
			name: "Namespace is WORD",
			want: &Token{
				Type:  WORD,
				Value: "namespace",
			},
		}, {
			name: "= is EQUAL",
			want: &Token{
				Type:  EQUALS,
				Value: "=",
			},
		},
		{
			name: "cooking is WORD",
			want: &Token{
				Type:  WORD,
				Value: "cooking",
			},
		},
		{
			name: "entity is WORD",
			want: &Token{
				Type:  WORD,
				Value: "entity",
			},
			skipnext: true,
		},
		{
			name: "{ is START",
			want: &Token{
				Type:  START,
				Value: "{",
			},
		},
		{
			name: "scope:kek is WORD",
			want: &Token{
				Type:  WORD,
				Value: "scope:character",
			},
			skipnext: true,
		},
		{
			name: "character.123 is WORD",
			want: &Token{
				Type:  WORD,
				Value: "character.123",
			},
			skipnext: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lexer.GetNextToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Lexer.GetNextToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
			if tt.posteffect != nil {
				tt.posteffect()
			}
			if tt.skipnext {
				lexer.GetNextToken()
			}
		})
	}
}

func TestLexer_Scan(t *testing.T) {
	tests := []struct {
		name    string
		lexer   *Lexer
		want    []*Token
		wantErr bool
	}{
		{
			name:  "elementary is tokenized correctly",
			lexer: New([]byte(elementary)),
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
			wantErr: false,
		},
		{
			name:  "elementary is tokenized correctly",
			lexer: New([]byte(scripted_trigger)),
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.lexer.Scan()
			if (err != nil) != tt.wantErr {
				t.Errorf("Lexer.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Stream, tt.want) {
				t.Errorf("Lexer.Scan() = %v, \nwant %v", got, tt.want)
			}
		})
	}
}
