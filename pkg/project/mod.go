package project

import (
	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/validator"
)

type ModFile struct {
	file             *files.FileEntry
	AST              *ast.AST
	Version          *tokens.Token   `json:"version"`
	Tags             []*tokens.Token `json:"tags"`
	Name             *tokens.Token   `json:"name"`
	Path             *tokens.Token   `json:"path"`
	ReplacePaths     []*tokens.Token `json:"replace_paths"`
	SupportedVersion *tokens.Token   `json:"supported_version"`
	Picture          *tokens.Token   `json:"picture"`
}

func NewModFile(AST *ast.AST, file_entry *files.FileEntry) *ModFile {
	block := AST.Block

	return &ModFile{
		file:             file_entry,
		AST:              AST,
		Version:          block.GetFieldValue("version"),
		Tags:             block.GetFieldList("tags"),
		Name:             block.GetFieldValue("name"),
		Path:             block.GetFieldValue("path"),
		ReplacePaths:     block.GetFieldsValues("replace_paths"),
		SupportedVersion: block.GetFieldValue("supported_version"),
		Picture:          block.GetFieldValue("picture"),
	}
}

func (m *ModFile) Validate() []*report.DiagnosticItem {
	diagnostics := make([]*report.DiagnosticItem, 0)

	fields := validator.NewBlockValidator(m.AST.Block)

	// Check for required fields
	fields.RequireField("version")
	fields.RequireField("name")
	fields.RequireField("path")

	// Check types of fields (if they exist)
	fields.ExpectString("version")
	fields.ExpectString("name")
	fields.ExpectString("path")

	// Optional fields: only check types if they are present
	fields.ExpectString("supported_version")
	fields.ExpectString("picture")

	diagnostics = append(diagnostics, fields.Errors()...)

	// validate token block
	tags := m.AST.Block.GetTokenBlock("tags")

	tag_validator := validator.NewTokenValidator(tags)
	tag_validator.ExpectAllTokensToBe(tokens.QUOTED_STRING)

	diagnostics = append(diagnostics, tag_validator.Errors()...)

	return diagnostics
}
