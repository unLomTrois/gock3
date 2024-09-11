package lexer

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNewTokenPatternMatcher(t *testing.T) {
	tpm := NewTokenPatternMatcher()

	if tpm == nil {
		t.Errorf("NewTokenPatternMatcher() returned nil")
	}

	if len(tpm.compiledRegexMap) == 0 {
		t.Errorf("NewTokenPatternMatcher() did not compile any regexes")
	}

	// Check if all expected token types are present
	expectedTokenTypes := []TokenType{
		COMMENT, WORD, STRING, NUMBER, BOOL, NEXTLINE,
		EQUALS, START, END, WHITESPACE, TAB, COMPARISON,
	}

	for _, tokenType := range expectedTokenTypes {
		if _, exists := tpm.compiledRegexMap[tokenType]; !exists {
			t.Errorf("NewTokenPatternMatcher() did not compile regex for TokenType %v", tokenType)
		}
	}
}

func TestTokenPatternMatcher_compileRegexes(t *testing.T) {
	tpm := &TokenPatternMatcher{
		compiledRegexMap: make(map[TokenType]*regexp.Regexp),
	}

	tpm.compileRegexes()

	if len(tpm.compiledRegexMap) == 0 {
		t.Errorf("compileRegexes() did not compile any regexes")
	}

	// Check if all expected token types are present
	expectedTokenTypes := []TokenType{
		COMMENT, WORD, STRING, NUMBER, BOOL, NEXTLINE,
		EQUALS, START, END, WHITESPACE, TAB, COMPARISON,
	}

	for _, tokenType := range expectedTokenTypes {
		if _, exists := tpm.compiledRegexMap[tokenType]; !exists {
			t.Errorf("compileRegexes() did not compile regex for TokenType %v", tokenType)
		}
	}
}

func TestTokenPatternMatcher_MatchToken(t *testing.T) {
	tpm := NewTokenPatternMatcher()

	tests := []struct {
		name      string
		tokenType TokenType
		text      []byte
		want      []byte
	}{
		{
			name:      "Match COMMENT token",
			tokenType: COMMENT,
			text:      []byte("# This is a comment\nAnd this is not!"),
			want:      []byte("# This is a comment"),
		},
		{
			name:      "Match WORD token",
			tokenType: WORD,
			text:      []byte("key = value"),
			want:      []byte("key"),
		},
		{
			name:      "Match WORD token with scope",
			tokenType: WORD,
			text:      []byte("scope:character = value"),
			want:      []byte("scope:character"),
		},
		{
			name:      "Match WORD token with dot notation",
			tokenType: WORD,
			text:      []byte("key.subkey = value"),
			want:      []byte("key.subkey"),
		},
		{
			name:      "Match STRING token",
			tokenType: STRING,
			text:      []byte(`"This is a string"\nAnd this is not`),
			want:      []byte(`"This is a string"`),
		},
		{
			name:      "Match NUMBER token - integer",
			tokenType: NUMBER,
			text:      []byte("123 not a number"),
			want:      []byte("123"),
		},
		{
			name:      "Match NUMBER token - float",
			tokenType: NUMBER,
			text:      []byte("123.45 not a number"),
			want:      []byte("123.45"),
		},
		{
			name:      "Match NUMBER token - negative",
			tokenType: NUMBER,
			text:      []byte("-123 not a number"),
			want:      []byte("-123"),
		},
		{
			name:      "Match BOOL token - yes",
			tokenType: BOOL,
			text:      []byte("yes no"),
			want:      []byte("yes"),
		},
		{
			name:      "Match BOOL token - no",
			tokenType: BOOL,
			text:      []byte("no yes"),
			want:      []byte("no"),
		},
		{
			name:      "Match NEXTLINE token",
			tokenType: NEXTLINE,
			text:      []byte("\n\nNext line"),
			want:      []byte("\n\n"),
		},
		{
			name:      "Match EQUALS token - single",
			tokenType: EQUALS,
			text:      []byte("= value"),
			want:      []byte("="),
		},
		{
			name:      "Match EQUALS token - double",
			tokenType: EQUALS,
			text:      []byte("== value"),
			want:      []byte("=="),
		},
		{
			name:      "Match START token",
			tokenType: START,
			text:      []byte("{ key = value }"),
			want:      []byte("{"),
		},
		{
			name:      "Match END token",
			tokenType: END,
			text:      []byte("} next"),
			want:      []byte("}"),
		},
		{
			name:      "Match WHITESPACE token",
			tokenType: WHITESPACE,
			text:      []byte("   next"),
			want:      []byte("   "),
		},
		{
			name:      "Match TAB token",
			tokenType: TAB,
			text:      []byte("\t\tnext"),
			want:      []byte("\t\t"),
		},
		{
			name:      "Match COMPARISON token - less than",
			tokenType: COMPARISON,
			text:      []byte("< 5"),
			want:      []byte("<"),
		},
		{
			name:      "Match COMPARISON token - greater than or equal",
			tokenType: COMPARISON,
			text:      []byte(">= 10"),
			want:      []byte(">="),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tpm.MatchToken(tt.tokenType, tt.text)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenPatternMatcher.MatchToken() = %s, want %s", got, tt.want)
			}
		})
	}
}
