// internal/app/lexer/lexer.go
package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
)

type Lexer struct {
	text      []byte
	cursor    int
	line      int
	regexpMap map[TokenType]*regexp.Regexp
}

// NormalizeText trims spaces and converts CRLF to LF
func NormalizeText(text []byte) []byte {
	text = bytes.TrimSpace(text)
	text = bytes.ReplaceAll(text, []byte("\r\n"), []byte("\n"))
	return regexp.MustCompile(`\n{3,}`).ReplaceAll(text, []byte("\n\n"))
}

// New creates a new Lexer instance
func New(text []byte) (*Lexer, error) {
	normalized := NormalizeText(text)

	if err := saveNormalizedText(normalized); err != nil {
		return nil, fmt.Errorf("failed to save normalized text: %w", err)
	}

	return &Lexer{
		text:      normalized,
		cursor:    0,
		line:      1,
		regexpMap: CompileRegexes(),
	}, nil
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
func (l *Lexer) Scan() (*TokenStream, error) {
	tokenStream := NewTokenStream()

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

func (l *Lexer) getNextToken() (*Token, error) {
	l.text = l.text[l.cursor:]

	for _, tokenType := range TokenCheckOrder {
		reg := l.regexpMap[tokenType]
		match := matchToken(reg, l.text)
		if match == nil {
			continue
		}

		l.cursor = len(match)

		switch tokenType {
		case WHITESPACE, TAB:
			return l.getNextToken()
		case NEXTLINE:
			l.line++
			return l.getNextToken()
		default:
			return &Token{
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
