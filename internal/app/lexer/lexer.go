package lexer

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
	"strconv"
)

type Lexer struct {
	Text   []byte
	Cursor int
	Line   int
}

// trims spaces and converts crlf to lf
func NormalizeText(text []byte) []byte {
	text = bytes.TrimSpace(text)
	text = bytes.ReplaceAll(text, []byte("\r\n"), []byte("\n"))
	// replace \n\n\n.. with \n\n, so it's only one line
	reg := regexp.MustCompile(`\n{3,}`)
	text = reg.ReplaceAll(text, []byte("\n\n"))

	return text
}

func New(text []byte) *Lexer {
	normalized := NormalizeText(text)

	new_file, _ := os.Create("./tmp/normalized.txt")
	defer new_file.Close()

	w := bufio.NewWriter(new_file)
	w.Write(normalized)
	w.Flush()

	return &Lexer{
		Text:   normalized,
		Cursor: 0,
		Line:   1,
	}
}

func (l *Lexer) hasMoreTokens() bool {
	return l.Cursor < len(l.Text)
}

func (l *Lexer) Scan() (*TokenStream, error) {

	tokenstream := NewTokenStream()

	for {
		if !l.hasMoreTokens() {
			break
		}

		token, err := l.GetNextToken()
		if err != nil {
			return nil, err
		}
		if token == nil {
			continue
		}

		tokenstream.Push(token)
	}

	return tokenstream, nil
}

func (l *Lexer) GetNextToken() (*Token, error) {

	l.Text = l.Text[l.Cursor:]

	for _, tokentype := range TokenCheckOrder {
		reg := regexp.MustCompile(TokenTypeToRegex[tokentype])
		match := l.match(reg, l.Text)
		l.Cursor = len(match)
		if match == nil {
			continue
		}
		if tokentype == WHITESPACE || tokentype == TAB {
			return l.GetNextToken()
		}
		if tokentype == NEXTLINE {
			// l.Line++
			return l.GetNextToken()
		}
		return &Token{
			Type:  tokentype,
			Value: match,
		}, nil
	}
	panic("[Lexer] Unexpected token: " + strconv.Quote(string(l.Text[0])))
}

func (l *Lexer) match(reg *regexp.Regexp, text []byte) []byte {
	match := reg.Find(text)
	if match == nil {
		return nil
	}
	return match
}

func (l *Lexer) GetContext(window int) string {
	if l.Cursor < len(l.Text) {
		end := l.Cursor + window
		if end > len(l.Text) {
			end = len(l.Text)
		}
		return string(l.Text[l.Cursor:end])
	}
	return ""
}
