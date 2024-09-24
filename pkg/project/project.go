package project

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/pdxfile"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

type Project struct {
	GameDir           string
	ModFileDescriptor string
	Diagnostics       []*report.DiagnosticItem
}

func NewProject(gameDir string, modFileDescriptor string) (*Project, error) {

	fmt.Println(gameDir, modFileDescriptor)

	// check game dir
	if _, err := os.Stat(gameDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("game directory %s does not exist", gameDir)
	}
	// check mod file
	if _, err := os.Stat(modFileDescriptor); os.IsNotExist(err) {
		return nil, fmt.Errorf("mod file %s does not exist", modFileDescriptor)
	}

	return &Project{
		GameDir:           gameDir,
		ModFileDescriptor: modFileDescriptor,
	}, nil
}

func (p *Project) Load() {

	p.LoadMod()

	p.Validate()

	// replace_paths := []string{}

	// mod_loader := files.NewModLoader(p.ModFileDescriptor, replace_paths)

	// fset := files.NewFileSet(p.GameDir, mod_loader)

	// err := fset.Scan()
}

func (p *Project) LoadMod() error {
	file_entry := files.NewFileEntry(p.ModFileDescriptor, files.FileKind(files.Mod))

	AST, err := pdxfile.ParseFile(file_entry)
	if err != nil {
		return err
	}

	mod := NewModFile(AST, file_entry)

	diagnostics := mod.Validate()
	p.Diagnostics = append(p.Diagnostics, diagnostics...)

	return nil
	// validate

	// mod_loader := files.NewModLoader(p.ModFileDescriptor, []string{})

	// fset := files.NewFileSet(p.GameDir, mod_loader)

	// err := fset.Scan()
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func (p *Project) Validate() []*report.DiagnosticItem {
	if len(p.Diagnostics) > 0 {
		for _, err := range p.Diagnostics {
			var c *color.Color
			switch err.Severity {
			case severity.Error:
				c = color.New(color.FgRed)
			case severity.Warning:
				c = color.New(color.FgYellow)
			case severity.Info:
				c = color.New(color.FgCyan)
			}
			filename, _ := err.Pointer.Loc.Filename()
			column := err.Pointer.Loc.Column
			line := err.Pointer.Loc.Line
			c.Println(fmt.Sprintf("[%s:%d:%d]: %s", filename, line, column, err.Msg))
		}
	}

	return p.Diagnostics
}
