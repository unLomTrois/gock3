package data

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/internal/app/pdxfile"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/validator"
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

func (traits *Traits) Load(fileEntries []*files.FileEntry) {
	traitFiles := traits.filterTraitFiles(fileEntries)

	log.Printf("Found %d files", len(traitFiles))
	if len(traitFiles) > 0 {
		log.Println(traitFiles[0].FullPath())
	}

	var problems []*report.DiagnosticItem

	for _, file := range traitFiles {
		ast := traits.loadFile(file)
		if ast == nil {
			continue
		}

		traitEntries, diagnostics := traits.parseTraits(ast.Block)
		traits.Traits = append(traits.Traits, traitEntries...)
		problems = append(problems, diagnostics...)
	}

	log.Printf("Found %d traits", len(traits.Traits))
	log.Printf("%d problems", len(problems))
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

func (traits *Traits) loadFile(fileEntry *files.FileEntry) *ast.AST {
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

	return traitEntries, problems
}

type Trait struct {
	key   *tokens.Token
	block *ast.FieldBlock
}

func NewTraitFromAST(key *tokens.Token, block *ast.FieldBlock) *Trait {
	return &Trait{
		key:   key,
		block: block,
	}
}

func (trait *Trait) Validate() []*report.DiagnosticItem {
	fields := validator.NewBlockValidator(trait.block)
	fields.ExpectNumber("minimum_age")
	fields.ExpectNumber("maximum_age")
	fields.ExpectNumber("intrigue")

	return fields.Errors()
}
