package main

import (
	// "ck3-parser/internal/app/linter"
	"ck3-parser/internal/app/parser"
	"encoding/json"
	"log"
	"os"
)

func main() {
	// Open file
	filepath := "data/0_elementary.txt"
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Parse file
	parser, err := parser.New(file)
	if err != nil {
		panic(err)
	}

	p, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	// Save parsed data
	err = SaveJSON(p, "parsed.json")
	if err != nil {
		panic(err)
	}

	// // Lint file
	// linter := linter.NewLinter(p.Filepath, p.Data)
	// linter.Lint()

	// lintedFilePath := "tmp/linted.txt"
	// err = SaveLintedData(linter, lintedFilePath)
	// if err != nil {
	// 	panic(err)
	// }
}

func SaveJSON(data interface{}, filename string) error {
	err := os.MkdirAll("tmp", 0755)
	if err != nil {
		return err
	}

	filepath := "tmp/" + filename
	file, err := os.Create(filepath)
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

	log.Println("Parsed data saved to tmp/parsed.json")

	return nil
}

// func SaveLintedData(linter *linter.Linter, filepath string) error {
// 	file, err := os.Create(filepath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	_, err = file.Write(linter.LintedData())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
