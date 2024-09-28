package utils

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/unLomTrois/gock3/pkg/cache"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

func (lex *Lexer) finalize() {

	diagnostics := lex.Errors()
	file_cache := cache.NewFileCache()

	for _, err := range diagnostics {
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

		c.Println(fmt.Sprintf("[%s:%d:%d]: %s, got %s", filename, line, column, err.Msg, err_line))
	}
}

func getErrorLine(fileCache *cache.FileCache, err *report.DiagnosticItem, column uint16) string {
	line_start := fileCache.GetLine(&err.Pointer.Loc)
	// fmt.Println(strconv.Quote(lineStart))

	// replace tabs to spaces, because loc sees \t as 4 symbols...
	// todo: do something
	spaced_line := strings.ReplaceAll(line_start, "\t", "    ")

	errorEndIndex := column + uint16(err.Pointer.Length) - 1
	return spaced_line[:errorEndIndex]
}
