package parser

import (
	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
)

// Parser represents the parser with its current state and error manager.
type Parser struct {
	tokenstream  *tokens.TokenStream
	currentToken *tokens.Token
	lookahead    *tokens.Token
	loc          *tokens.Loc
	*report.ErrorManager
}

// New creates a new Parser instance.
func New(tokenstream *tokens.TokenStream) *Parser {
	p := &Parser{
		tokenstream:  tokenstream,
		ErrorManager: report.NewErrorManager(),
	}
	p.currentToken = p.tokenstream.Next()
	p.lookahead = p.tokenstream.Next()
	if p.currentToken != nil {
		p.loc = &p.currentToken.Loc
	}
	return p
}

// Parse processes the token stream and returns the AST along with any diagnostic errors.
func Parse(token_stream *tokens.TokenStream) (*ast.FileBlock, []*report.DiagnosticItem) {
	p := New(token_stream)
	fileBlock := p.fileBlock()
	return fileBlock, p.Errors()
}

// nextToken advances the currentToken and lookahead tokens.
func (p *Parser) nextToken() {
	p.currentToken = p.lookahead
	p.lookahead = p.tokenstream.Next()
	if p.currentToken != nil {
		p.loc = &p.currentToken.Loc
	}
}
