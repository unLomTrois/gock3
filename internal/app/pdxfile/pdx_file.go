package pdxfile

import (
	"ck3-parser/internal/app/files"
	"ck3-parser/internal/app/lexer"
	"ck3-parser/internal/app/parser"
	"ck3-parser/internal/app/utils"
	"fmt"
	"os"
)

func ParseFile(entry *files.FileEntry) ([]*parser.Node, error) {
	content, err := os.ReadFile(entry.FullPath())
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	token_stream, err := lexer.Scan(content)
	if err != nil {
		return nil, err
	}

	utils.SaveJSON(token_stream.Stream, "tmp/token_stream.json")

	p := parser.New(token_stream)
	parse_tree := p.Parse()
	// todo: err here

	return parse_tree, nil
}
