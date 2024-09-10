package main

import (
	"ck3-parser/internal/app/lexer"
	"ck3-parser/internal/app/linter"
	"ck3-parser/internal/app/parser"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	inputFilePath = "data/0_elementary.txt"
	outputDir     = "tmp"
	parseTreeFile = "parsetree.json"
	lintedFile    = "linted.txt"
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

	content, err := readFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	tokenStream, err := scanContent(content)
	if err != nil {
		return fmt.Errorf("scanning content: %w", err)
	}

	parseTrees, err := parseTokens(tokenStream)
	if err != nil {
		return fmt.Errorf("parsing tokens: %w", err)
	}

	if err := saveJSON(parseTrees, parseTreeFile); err != nil {
		return fmt.Errorf("saving parse tree: %w", err)
	}

	if err := lintAndSave(parseTrees); err != nil {
		return fmt.Errorf("linting and saving: %w", err)
	}

	return nil
}

func readFile(path string) ([]byte, error) {
	start := time.Now()
	defer func() {
		log.Printf("File read time: %s", time.Since(start))
	}()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func scanContent(content []byte) (*lexer.TokenStream, error) {
	start := time.Now()
	defer func() {
		log.Printf("Scan time: %s", time.Since(start))
	}()

	l := lexer.NewLexer(content)
	tokenStream, err := l.Scan()
	if err != nil {
		return nil, fmt.Errorf("failed to scan content: %w", err)
	}

	return tokenStream, nil
}

func parseTokens(tokens *lexer.TokenStream) ([]*parser.Node, error) {
	start := time.Now()
	defer func() {
		log.Printf("Parse time: %s", time.Since(start))
	}()

	p := parser.New(tokens)
	return p.Parse(), nil
}

func saveJSON(data interface{}, filename string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(outputDir, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	if err := enc.Encode(data); err != nil {
		return err
	}

	log.Printf("Saved JSON to %s", filepath.Join(outputDir, filename))
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
