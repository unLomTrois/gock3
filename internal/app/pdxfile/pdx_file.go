package pdxfile

import (
	"fmt"
	"os"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer"
	"github.com/unLomTrois/gock3/internal/app/parser"
	"github.com/unLomTrois/gock3/internal/app/utils"
)

func ParseFile(entry *files.FileEntry) (*parser.FileBlock, error) {
	content, err := os.ReadFile(entry.FullPath())
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	token_stream, err := lexer.Scan(entry, content)
	if err != nil {
		return nil, err
	}

	utils.SaveJSON(token_stream.Stream, "tmp/token_stream.json")

	parse_tree := parser.Parse(token_stream)
	// todo: err here

	return parse_tree, nil
}
