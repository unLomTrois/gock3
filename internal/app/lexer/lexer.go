// internal/app/lexer/lexer.go
package lexer

import (
	"bufio"
	"bytes"
	"ck3-parser/internal/app/tokens"
	"fmt"
	"os"
)

type Lexer struct {
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
func NewLexer(text []byte) *Lexer {
	return &Lexer{
		text:           NormalizeText(text),
		cursor:         0,
		line:           1,
		patternMatcher: NewTokenPatternMatcher(),
	}
}

func saveNormalizedText(text []byte) error {
	file, err := os.Create("./tmp/normalized.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(text); err != nil {
		return err
	}
	return writer.Flush()
}

func (l *Lexer) hasMoreTokens() bool {
	return l.cursor < len(l.text)
}

// Scan tokenizes the entire input text
func (l *Lexer) Scan() (*tokens.TokenStream, error) {
	tokenStream := tokens.NewTokenStream()

	for l.hasMoreTokens() {
		token, err := l.getNextToken()
		if err != nil {
			return nil, fmt.Errorf("error scanning tokens: %w", err)
		}
		if token != nil {
			tokenStream.Push(token)
		}
	}

	return tokenStream, nil
}

func (l *Lexer) getNextToken() (*tokens.Token, error) {
	l.text = l.text[l.cursor:]

	for _, tokenType := range tokens.TokenCheckOrder {
		match := l.patternMatcher.MatchToken(tokenType, l.text)
		if match == nil {
			continue
		}

		l.cursor = len(match)

		switch tokenType {
		case tokens.WHITESPACE, tokens.TAB:
			return l.getNextToken()
		case tokens.NEXTLINE:
			l.line++
			return l.getNextToken()
		default:
			return &tokens.Token{
				Type:  tokenType,
				Value: string(match),
			}, nil
		}
	}

	return nil, fmt.Errorf("unexpected token at position: line %d, col %d: %q", l.line, l.cursor, string(l.text[0]))
}

// GetContext returns a window of characters around the current cursor position
func (l *Lexer) GetContext(window int) string {
	if l.cursor >= len(l.text) {
		return ""
	}
	end := l.cursor + window
	if end > len(l.text) {
		end = len(l.text)
	}
	return string(l.text[l.cursor:end])
}
