package main

import (
	"ck3-parser/internal/app/files"
	"ck3-parser/internal/app/linter"
	"ck3-parser/internal/app/parser"
	"ck3-parser/internal/app/pdxfile"
	"ck3-parser/internal/app/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	inputFilePath   = "data/3_traits.txt"
	outputDir       = "tmp"
	tokenStreamFile = "token_stream.json"
	parseTreeFile   = "parsetree.json"
	lintedFile      = "linted.txt"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	start := time.Now()
	defer func() {
		log.Printf("Total execution time: %s", time.Since(start))
	}()

	path := inputFilePath
	fullpath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}
	file_entry := files.NewFileEntry(path, fullpath, files.FileKind(files.Mod))

	parseTrees, err := pdxfile.ParseFile(file_entry)
	if err != nil {
		return fmt.Errorf("parsing tokens: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	if err := utils.SaveJSON(parseTrees, filepath.Join(outputDir, parseTreeFile)); err != nil {
		return fmt.Errorf("saving parse tree: %w", err)
	}

	if err := lintAndSave(parseTrees); err != nil {
		return fmt.Errorf("linting and saving: %w", err)
	}

	return nil
}

func lintAndSave(parseTrees []*parser.Node) error {
	config := linter.LintConfig{
		IntendStyle:            linter.TABS,
		IntendSize:             4,
		TrimTrailingWhitespace: true,
		InsertFinalNewline:     true,
		CharSet:                "utf-8-bom",
		EndOfLine:              []byte("\r\n"),
	}

	l := linter.New(parseTrees, config)
	l.Lint()

	if err := l.Save(filepath.Join(outputDir, lintedFile)); err != nil {
		return err
	}

	log.Printf("Linted file saved to %s", filepath.Join(outputDir, lintedFile))
	return nil
}
