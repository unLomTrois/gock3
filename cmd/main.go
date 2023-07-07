package main

import (
	// "ck3-parser/internal/app/linter"

	"ck3-parser/internal/app/lexer"
	"ck3-parser/internal/app/parser"
	"encoding/json"
	"io"
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

	filecontent, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	lexer := lexer.New(filecontent)
	tokenstream, err := lexer.Scan()
	if err != nil {
		panic(err)
	}

	err = SaveJSON(tokenstream, "tokenstream.json")
	if err != nil {
		panic(err)
	} else {
		log.Println("Parsed data saved to tmp/tokenstream.json")
	}

	parser := parser.New(tokenstream)
	parsetree := parser.Parse()
	// parser.Parse()

	SaveJSON(parsetree, "parsetree.json")
	log.Println("Parsed data saved to tmp/parsetree.json")

	// Lint file
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
