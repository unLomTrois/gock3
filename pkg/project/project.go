package project

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/pdxfile"
	"github.com/unLomTrois/gock3/pkg/cache"
	"github.com/unLomTrois/gock3/pkg/data"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

type Project struct {
	VanillaDir        string
	ModFileDescriptor string
	Diagnostics       []*report.DiagnosticItem
	Common            *data.Common
}

func NewProject(vanillaDir string, modFileDescriptor string) (*Project, error) {
	// check game dir
	if _, err := os.Stat(vanillaDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("game directory %s does not exist", vanillaDir)
	}
	// check mod file
	if _, err := os.Stat(modFileDescriptor); os.IsNotExist(err) {
		return nil, fmt.Errorf("mod file %s does not exist", modFileDescriptor)
	}

	return &Project{
		VanillaDir:        vanillaDir,
		ModFileDescriptor: modFileDescriptor,
		Diagnostics:       []*report.DiagnosticItem{},
		Common:            data.NewCommon(),
	}, nil
}

func (project *Project) Load() {
	mod := project.LoadMod()

	replacePaths := make([]string, len(mod.ReplacePaths))
	for i, token := range mod.ReplacePaths {
		replacePaths[i] = token.Value
	}

	log.Println("mod path", mod.Path)

	modLoader := files.NewModLoader(mod.Path.Value, replacePaths)
	fset := files.NewFileSet(project.VanillaDir, modLoader)

	fset.Scan(project.VanillaDir)

	project.Common.Load(fset)

	project.Validate()
}

func (p *Project) LoadMod() *ModFile {
	file_entry := files.NewFileEntry(p.ModFileDescriptor, files.FileKind(files.Mod))

	AST, err := pdxfile.ParseFile(file_entry)
	if err != nil {
		return nil
	}

	mod := NewModFile(AST, file_entry)

	diagnostics := mod.Validate()
	p.Diagnostics = append(p.Diagnostics, diagnostics...)

	return mod
}

func (p *Project) Validate() []*report.DiagnosticItem {
	file_cache := cache.NewFileCache()

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

			err_line := getErrorLine(file_cache, err, column)

			if err.Pointer.Loc.Line == 1 && err.Pointer.Loc.Column == 1 {
				c.Println(fmt.Sprintf("[%s:%d:%d]: %s", filename, line, column, err.Msg))

				continue
			}

			c.Println(fmt.Sprintf("[%s:%d:%d]: %s, got %s", filename, line, column, err.Msg, err_line))
		}
	}

	return p.Diagnostics
}

func getErrorLine(fileCache *cache.FileCache, err *report.DiagnosticItem, column uint16) string {
	line_start := fileCache.GetLine(&err.Pointer.Loc)
	// fmt.Println(strconv.Quote(lineStart))

	// replace tabs to spaces, because loc sees \t as 4 symbols...
	// todo: do something
	spaced_line := strings.ReplaceAll(line_start, "\t", "    ")

	errorEndIndex := column + uint16(err.Pointer.Length) - 1
	return spaced_line[:errorEndIndex]
}
