package report

import (
	"fmt"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

type DiagnosticItem struct {
	Severity severity.Severity
	Pointer  *DiagnosticPointer
	Msg      string
}

type DiagnosticPointer struct {
	Loc    *files.Loc
	Length int
}

func (d *DiagnosticItem) Error() string {
	return fmt.Sprintf("%s: %s", d.Severity, d.Msg)
}

func NewDiagnosticItem(severity severity.Severity, msg string, pointer *DiagnosticPointer) *DiagnosticItem {
	return &DiagnosticItem{
		Severity: severity,
		Msg:      msg,
		Pointer:  pointer,
	}
}

func FromToken(token *tokens.Token, severity severity.Severity, msg string) *DiagnosticItem {
	return &DiagnosticItem{
		Severity: severity,
		Msg:      msg,
		Pointer: &DiagnosticPointer{
			Loc:    token.Loc,
			Length: len(token.Value),
		},
	}
}

func FromFile(file *files.FileEntry, severity severity.Severity, msg string) *DiagnosticItem {
	loc := files.LocFromFileEntry(file)

	return &DiagnosticItem{
		Severity: severity,
		Msg:      msg,
		Pointer: &DiagnosticPointer{
			Loc:    loc,
			Length: 0,
		},
	}
}
