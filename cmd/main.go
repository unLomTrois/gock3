package main

import (
	"ck3-parser/internal/app/files"
	"ck3-parser/internal/app/pdxfile"
	"ck3-parser/internal/app/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	inputFilePath   = "data/5_event.txt"
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

	vanilla_root := ""
	mod_root := ""
	replace_paths := []string{}

	mod_loader := files.NewModLoader(mod_root, replace_paths)

	fset := files.NewFileSet(vanilla_root, mod_loader)

	traits_dir := "C:/Users/vadim/Documents/Paradox Interactive/Crusader Kings III/mod/T4N-CK3/T4N/common/traits"

	err := fset.Scan(traits_dir)
	if err != nil {
		return fmt.Errorf("scanning files: %w", err)
	}

	path := inputFilePath
	fullpath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("getting absolute path: %w", err)
	}
	file_entry := files.NewFileEntry(fullpath, files.FileKind(files.Mod))

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

	return nil
}
