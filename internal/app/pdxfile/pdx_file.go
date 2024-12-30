package pdxfile

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer"
	"github.com/unLomTrois/gock3/internal/app/parser"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/internal/app/utils"
	"github.com/unLomTrois/gock3/pkg/cache"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

func ParseFile(entry *files.FileEntry) (*ast.AST, error) {
	content, err := utils.ReadFileWithUTF8BOM(entry.FullPath())
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var errs []*report.DiagnosticItem

	token_stream, lexer_errs := lexer.Scan(entry, content)

	// utils.SaveJSON(token_stream, "tokenstream.json")

	errs = append(errs, lexer_errs...)

	file_block, parser_errs := parser.Parse(token_stream)
	errs = append(errs, parser_errs...)

	ast := &ast.AST{
		Filename: entry.FileName(),
		Fullpath: entry.FullPath(),
		Block:    file_block,
	}

	finalize(errs)

	return ast, nil
}

func finalize(errs []*report.DiagnosticItem) {

	file_cache := cache.NewFileCache()

	for _, err := range errs {
		var c *color.Color
		switch err.Severity {
		case severity.Error:
			c = color.New(color.FgRed)
		case severity.Warning:
			c = color.New(color.FgYellow)
		case severity.Info:
			c = color.New(color.FgCyan)
		case severity.Critical:
			c = color.New(color.FgHiMagenta)
		}
		filename, _ := err.Pointer.Loc.Filename()
		column := err.Pointer.Loc.Column
		line := err.Pointer.Loc.Line

		err_line := getErrorLine(file_cache, err, column)

		if err.Pointer.Loc.Line == 1 && err.Pointer.Loc.Column == 1 {
			c.Println(fmt.Sprintf("[%s:%d:%d]: %s", filename, line, column, err.Msg))

			continue
		}

		c.Println(fmt.Sprintf("[%s:%d:%d]: %s, got %s", filename, line, column, err.Msg, strconv.Quote(err_line)))
	}
}

func getErrorLine(fileCache *cache.FileCache, err *report.DiagnosticItem, column uint16) string {
	lineStart := fileCache.GetLine(&err.Pointer.Loc)

	// replace tabs to spaces, because loc sees \t as 4 symbols...
	// todo: do something
	spacedLine := strings.ReplaceAll(lineStart, "\t", "    ")

	errorEndIndex := column + uint16(err.Pointer.Length) - 1

	return spacedLine[:errorEndIndex]
}
