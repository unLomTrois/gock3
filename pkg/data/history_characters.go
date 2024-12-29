package data

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/internal/app/pdxfile"
	"github.com/unLomTrois/gock3/pkg/report"
)

type HistoryCharacters struct {
	Characters []*HistoryCharacter
	ast        *ast.AST
}

func NewHistoryCharacters() *HistoryCharacters {
	return &HistoryCharacters{
		Characters: make([]*HistoryCharacter, 0),
		ast:        nil,
	}
}

// Folder returns the folder path for traits, using the correct
// path separator for the operating system.
func (t *HistoryCharacters) Folder() string {
	return filepath.Join("history", "characters")
}

func (hc *HistoryCharacters) Load(fileEntries []*files.FileEntry) []*HistoryCharacter {
	files := hc.filterFiles(fileEntries)

	log.Printf("Found %d character files", len(files))

	var problems []*report.DiagnosticItem

	for _, file := range files {
		ast := hc.parseFile(file)
		if ast == nil {
			continue
		}

		entries, diagnostics := hc.parse(ast.Block)
		hc.Characters = append(hc.Characters, entries...)
		problems = append(problems, diagnostics...)
	}

	log.Printf("Found %d characters", len(hc.Characters))
	log.Printf("%d problems", len(problems))

	return hc.Characters
}

func (hc *HistoryCharacters) filterFiles(fileEntries []*files.FileEntry) []*files.FileEntry {
	traitFiles := make([]*files.FileEntry, 0, len(fileEntries))
	for _, fileEntry := range fileEntries {
		if strings.Contains(fileEntry.FullPath(), hc.Folder()) {
			traitFiles = append(traitFiles, fileEntry)
		}
	}
	return traitFiles
}

func (hc *HistoryCharacters) parseFile(fileEntry *files.FileEntry) *ast.AST {
	ast, err := pdxfile.ParseFile(fileEntry)
	if err != nil {
		log.Printf("Failed to parse file %s: %v", fileEntry.FullPath(), err)
		return nil
	}
	return ast
}

func (traits *HistoryCharacters) parse(block *ast.FieldBlock) ([]*HistoryCharacter, []*report.DiagnosticItem) {
	var entities []*HistoryCharacter
	var problems []*report.DiagnosticItem

	for _, field := range block.Values {
		// Skip variables
		if strings.Contains(field.Key.Value, "@") {
			continue
		}

		key := field.Key
		block, ok := field.Value.(*ast.FieldBlock)
		if !ok {
			continue
		}

		character := NewHistoryCharacter(key, block)
		problems = append(problems, character.Validate()...)
		entities = append(entities, character)
	}

	// for _, trait := range traitEntries {
	// 	fp, err := trait.key.Loc.Fullpath()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	log.Println(trait.name, fp)
	// }

	return entities, problems
}
