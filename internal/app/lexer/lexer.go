// internal/app/lexer/lexer.go
package lexer

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
)

type Lexer struct {
	entry          *files.FileEntry
	text           []byte
	cursor         int
	line           int
	column         int
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
		column:         1,
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

	for lex.hasMoreTokens() {
		token, err := lex.getNextToken()
		if err != nil {
			return nil, fmt.Errorf("error scanning tokens: %w", err)
		}
		if token != nil {
			tokenStream.Push(token)
		}
	}

	return tokenStream, nil
}

func (lex *Lexer) remainder() []byte {
	return lex.text[lex.cursor:]
}

func (lex *Lexer) getNextToken() (*tokens.Token, error) {
	if !lex.hasMoreTokens() {
		return nil, nil
	}

	remainder := lex.remainder()

	var matchedToken []byte
	var matchedTokenType tokens.TokenType

	// Keep track of the initial line and column for the token
	startLine := lex.line
	startColumn := lex.column

	// Try to match tokens in the specified order
	for _, tokenType := range tokens.TokenCheckOrder {
		match := lex.patternMatcher.MatchToken(tokenType, remainder)
		if match == nil {
			continue
		}

		// Accept the first match
		matchedToken = match
		matchedTokenType = tokenType
		break
	}

	if matchedToken != nil {
		tokenValue := string(matchedToken)
		lex.cursor += len(matchedToken)

		switch matchedTokenType {
		case tokens.TAB:
			// Consider tab width as 4 spaces
			lex.column += 4 // Already added 1 in len(longestMatch)
			return nil, nil
		case tokens.NEXTLINE:
			lex.line++
			lex.column = 1
			return nil, nil
		case tokens.WHITESPACE:
			// Ignore whitespace
			return nil, nil
		case tokens.COMMENT:
			// Ignore comments
			return nil, nil
		default:
			lex.column += len(matchedToken)
			loc := tokens.LocFromFileEntry(lex.entry)
			loc.Line = uint32(startLine)
			loc.Column = uint16(startColumn)
			return tokens.New(tokenValue, matchedTokenType, *loc), nil
		}
	}

	// Handle unexpected characters
	unexpectedChar := remainder[0]
	errMsg := fmt.Sprintf("unexpected token at line %d, column %d: '%c'", lex.line, lex.column, unexpectedChar)

	// Advance cursor to prevent infinite loop
	lex.cursor++
	lex.column++

	return nil, errors.New(errMsg)
}
