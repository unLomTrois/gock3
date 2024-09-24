package pdxfile

import (
	"fmt"
	"os"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer"
	"github.com/unLomTrois/gock3/internal/app/parser"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
)

func ParseFile(entry *files.FileEntry) (*ast.AST, error) {
	content, err := os.ReadFile(entry.FullPath())
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	token_stream, err := lexer.Scan(entry, content)
	if err != nil {
		return nil, err
	}

	file_block := parser.Parse(token_stream)
	// todo: err here

	ast := &ast.AST{
		Filename: entry.FileName(),
		Fullpath: entry.FullPath(),
		Block:    file_block,
	}

	return ast, nil
}
