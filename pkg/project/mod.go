package project

import (
	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
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

	// check required fields
	if m.Version == nil {
		diagnostics = append(diagnostics, report.FromFile(
			m.file,
			severity.Error,
			"version is required",
		))
	}
	if m.Name == nil {
		diagnostics = append(diagnostics, report.FromFile(
			m.file,
			severity.Error,
			"name is required",
		))
	}
	if m.Path == nil {
		diagnostics = append(diagnostics, report.FromFile(
			m.file,
			severity.Error,
			"path is required",
		))
	}

	return diagnostics
}
