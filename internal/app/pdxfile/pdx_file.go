package pdxfile

import (
	"fmt"
	"os"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer"
	"github.com/unLomTrois/gock3/internal/app/parser"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
)

func ParseFile(entry *files.FileEntry) (*ast.AST, error) {
	content, err := os.ReadFile(entry.FullPath())
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var errs []*report.DiagnosticItem

	token_stream, lexer_errs := lexer.Scan(entry, content)

	errs = append(errs, lexer_errs...)

	file_block, parser_errs := parser.Parse(token_stream)
	errs = append(errs, parser_errs...)

	ast := &ast.AST{
		Filename: entry.FileName(),
		Fullpath: entry.FullPath(),
		Block:    file_block,
	}

	return ast, nil
}
