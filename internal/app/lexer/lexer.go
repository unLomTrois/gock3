// internal/app/lexer/lexer.go
package lexer

import (
	"bytes"
	"fmt"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

type Lexer struct {
	entry          *files.FileEntry
	text           []byte
	cursor         int
	line           int
	column         int
	patternMatcher *TokenPatternMatcher
	*report.ErrorManager
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
		ErrorManager:   report.NewErrorManager(),
	}
}

// NormalizeText trims spaces and converts CRLF to LF
func NormalizeText(text []byte) []byte {
	text = bytes.TrimSpace(text)
	text = bytes.ReplaceAll(text, []byte("\r\n"), []byte("\n"))
	return text
}

func (lex *Lexer) hasMoreTokens() bool {
	return lex.cursor < len(lex.text)
}

// Scan tokenizes the entire input text
func Scan(entry *files.FileEntry, text []byte) (*tokens.TokenStream, []*report.DiagnosticItem) {
	lex := NewLexer(entry, text)

	tokenStream := tokens.NewTokenStream()

	for lex.hasMoreTokens() {
		token := lex.getNextToken()
		if token == nil {
			continue
		}
		tokenStream.Push(token)
	}

	return tokenStream, lex.Errors()
}

func (lex *Lexer) remainder() []byte {
	return lex.text[lex.cursor:]
}

func (lex *Lexer) getNextToken() *tokens.Token {
	if !lex.hasMoreTokens() {
		return nil
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
			return nil
		case tokens.NEXTLINE:
			lex.line++
			lex.column = 1
			return nil
		case tokens.WHITESPACE:
			// Ignore whitespace
			return nil
		case tokens.COMMENT:
			// Ignore comments
			return nil
		default:
			lex.column += len(matchedToken)
			loc := tokens.LocFromFileEntry(lex.entry)
			loc.Line = uint32(startLine)
			loc.Column = uint16(startColumn)
			return tokens.New(tokenValue, matchedTokenType, *loc)
		}
	}

	unexpectedChar := remainder[0]
	loc := tokens.LocFromFileEntry(lex.entry)
	loc.Line = uint32(lex.line)
	loc.Column = uint16(lex.column)
	err := report.FromLoc(*loc, severity.Critical, fmt.Sprintf("unexpected token '%c'", unexpectedChar))
	lex.AddError(err)

	// Advance cursor to prevent infinite loop
	lex.cursor++
	lex.column++

	return nil
}
