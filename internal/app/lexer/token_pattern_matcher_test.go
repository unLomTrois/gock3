package lexer

import (
	"reflect"
	"testing"
)

func TestNewTokenPatternMatcher(t *testing.T) {
	tests := []struct {
		name string
		want *TokenPatternMatcher
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenPatternMatcher(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenPatternMatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenPatternMatcher_compileRegexes(t *testing.T) {
	tests := []struct {
		name string
		tpm  *TokenPatternMatcher
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tpm.compileRegexes()
		})
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
			name:      "Match SCRIPT token",
			tokenType: SCRIPT,
			text:      []byte("scripted_trigger"),
			want:      []byte("scripted_trigger"),
		},
		{
			name:      "Match WORD token",
			tokenType: WORD,
			text:      []byte("key = value"),
			want:      []byte("key"),
		},
		{
			name:      "Match hard WORD token",
			tokenType: WORD,
			text:      []byte("key_1.2 = value"),
			want:      []byte("key_1.2"),
		},
		{
			name:      "Match STRING token",
			tokenType: STRING,
			text:      []byte(`"This is a string"\nAnd this is not`),
			want:      []byte(`"This is a string"`),
		},
		{
			name:      "yes BOOL token",
			tokenType: BOOL,
			text:      []byte("yes\nno"),
			want:      []byte("yes"),
		},
		{
			name:      "no BOOL token",
			tokenType: BOOL,
			text:      []byte("no\nyes"),
			want:      []byte("no"),
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
