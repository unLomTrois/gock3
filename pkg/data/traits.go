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

type Traits struct {
	Traits []*Trait
	ast    *ast.AST
}

func NewTraits() *Traits {
	return &Traits{
		Traits: []*Trait{},
		ast:    &ast.AST{},
	}
}

// Folder returns the folder path for traits, using the correct
// path separator for the operating system.
func (t *Traits) Folder() string {
	return filepath.Join("common", "traits")
}

func (traits *Traits) Load(fileEntries []*files.FileEntry) []*Trait {
	traitFiles := traits.filterTraitFiles(fileEntries)

	log.Printf("Found %d trait files", len(traitFiles))

	var problems []*report.DiagnosticItem

	for _, file := range traitFiles {
		ast := traits.parseFile(file)
		if ast == nil {
			continue
		}

		traitEntries, diagnostics := traits.parseTraits(ast.Block)
		traits.Traits = append(traits.Traits, traitEntries...)
		problems = append(problems, diagnostics...)
	}

	log.Printf("Found %d traits", len(traits.Traits))
	log.Printf("%d problems", len(problems))

	return traits.Traits
}

func (traits *Traits) filterTraitFiles(fileEntries []*files.FileEntry) []*files.FileEntry {
	traitFiles := make([]*files.FileEntry, 0, len(fileEntries))
	for _, fileEntry := range fileEntries {
		if strings.Contains(fileEntry.FullPath(), traits.Folder()) {
			traitFiles = append(traitFiles, fileEntry)
		}
	}
	return traitFiles
}

func (traits *Traits) parseFile(fileEntry *files.FileEntry) *ast.AST {
	ast, err := pdxfile.ParseFile(fileEntry)
	if err != nil {
		log.Printf("Failed to parse file %s: %v", fileEntry.FullPath(), err)
		return nil
	}
	return ast
}

func (traits *Traits) parseTraits(block *ast.FieldBlock) ([]*Trait, []*report.DiagnosticItem) {
	var traitEntries []*Trait
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

		trait := NewTraitFromAST(key, block)
		problems = append(problems, trait.Validate()...)
		traitEntries = append(traitEntries, trait)
	}

	// for _, trait := range traitEntries {
	// 	fp, err := trait.key.Loc.Fullpath()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	log.Println(trait.name, fp)
	// }

	return traitEntries, problems
}
