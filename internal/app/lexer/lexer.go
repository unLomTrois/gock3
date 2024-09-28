// internal/app/lexer/lexer.go
package lexer

import (
	"bytes"
	"fmt"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
)

type Lexer struct {
	entry          *files.FileEntry
	text           []byte
	cursor         int
	line           int
	patternMatcher *TokenPatternMatcher
}

// NormalizeText trims spaces and converts CRLF to LF
func NormalizeText(text []byte) []byte {
	text = bytes.TrimSpace(text)
	text = bytes.ReplaceAll(text, []byte("\r\n"), []byte("\n"))
	return text
}

// NewLexer creates a new Lexer instance
func NewLexer(entry *files.FileEntry, text []byte) *Lexer {
	return &Lexer{
		entry:          entry,
		text:           NormalizeText(text),
		cursor:         0,
		line:           1,
		patternMatcher: NewTokenPatternMatcher(),
	}
}

func (lex *Lexer) hasMoreTokens() bool {
	return lex.cursor < len(lex.text)
}

// Scan tokenizes the entire input text
func Scan(entry *files.FileEntry, text []byte) (*tokens.TokenStream, error) {
	lex := NewLexer(entry, text)

	tokenStream := tokens.NewTokenStream()

	loc := tokens.LocFromFileEntry(entry)

	for lex.hasMoreTokens() {
		token, err := lex.getNextToken(loc)
		if err != nil {
			return nil, fmt.Errorf("error scanning tokens: %w", err)
		}
		if token != nil {
			tokenStream.Push(token)
			loc.Column += uint16(len(token.Value))
		}
	}

	return tokenStream, nil
}

func (lex *Lexer) remainder() []byte {
	return lex.text[lex.cursor:]
}

func (lex *Lexer) getNextToken(loc *tokens.Loc) (*tokens.Token, error) {
	remainder := lex.remainder()

	for _, tokenType := range tokens.TokenCheckOrder {
		match := lex.patternMatcher.MatchToken(tokenType, remainder)
		if match == nil {
			continue
		}

		lex.cursor += len(match)

		switch tokenType {
		case tokens.TAB:
			// TODO: Consider using a tab width from users editor settings
			loc.Column += 4
			return nil, nil
		case tokens.NEXTLINE:
			loc.Line += 1
			loc.Column = 1
			lex.line++
			return nil, nil
		case tokens.WHITESPACE:
			loc.Column += 1
			return nil, nil
		default:
			return tokens.New(string(match), tokenType, *loc), nil
		}
	}

	return nil, fmt.Errorf("unexpected token at position: line %d, col %d: %q", loc.Line, loc.Column, string(remainder[0]))
}
