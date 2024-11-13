package data

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/internal/app/pdxfile"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/validator"
)

type Traits struct {
	Traits []*Trait
}

func NewTraits() *Traits {
	return &Traits{
		Traits: []*Trait{},
	}
}

// Folder returns the folder path for traits, using the correct
// path separator for the operating system.
// On Windows, it returns "common\\traits", and on Linux, "common/traits".
func (t *Traits) Folder() string {
	return filepath.Join("common", "traits")
}

func (traits *Traits) Load(fileEntries []*files.FileEntry) {
	var traitFiles []*files.FileEntry

	for _, fileEntry := range fileEntries {
		fullpath := fileEntry.FullPath()
		if strings.Contains(fullpath, traits.Folder()) {
			traitFiles = append(traitFiles, fileEntry)
		}
	}

	log.Println("Found", len(traitFiles), "files")
	log.Println(traitFiles[0].FullPath())

	for _, entry := range traitFiles {
		ast := traits.LoadFile(entry)
		trait := NewTraitFromAST(ast)
		trait.Validate()
	}
}

func (traits *Traits) LoadFile(fileEntry *files.FileEntry) *ast.AST {
	AST, err := pdxfile.ParseFile(fileEntry)
	if err != nil {
		return nil
	}
	return AST
}

type Trait struct {
	ast *ast.AST
}

func NewTraitFromAST(ast *ast.AST) *Trait {
	return &Trait{
		ast: ast,
	}
}

func (trait *Trait) Validate() []*report.DiagnosticItem {
	diagnostics := make([]*report.DiagnosticItem, 0)

	fields := validator.NewBlockValidator(trait.ast.Block)

	// Check for required fields
	// fields.RequireField("version")
	// fields.RequireField("name")
	// fields.RequireField("path")

	// // Check types of fields (if they exist)
	// fields.ExpectString("version")
	// fields.ExpectString("name")
	// fields.ExpectString("path")

	// // Optional fields: only check types if they are present
	// fields.ExpectString("supported_version")
	// fields.ExpectString("picture")

	diagnostics = append(diagnostics, fields.Errors()...)

	// // validate token block
	// tags := trait.ast.Block.GetTokenBlock("tags")

	// tag_validator := validator.NewTokenValidator(tags)
	// tag_validator.ExpectAllTokensToBe(tokens.QUOTED_STRING)

	// diagnostics = append(diagnostics, tag_validator.Errors()...)

	return diagnostics

}
