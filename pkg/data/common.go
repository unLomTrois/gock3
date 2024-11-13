package data

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/unLomTrois/gock3/internal/app/files"
)

type Common struct {
	Traits *Traits
}

func NewCommon() *Common {
	return &Common{
		Traits: NewTraits(),
	}
}

// Folder returns the folder path for common, using the correct
// path separator for the operating system.
// On Windows, it returns "game\\common", and on Linux, "game/common".
func (c *Common) Folder() string {
	return filepath.Join("game", "common")
}

func (common *Common) Load(fset *files.FileSet) {
	var commonFiles []*files.FileEntry

	for _, fileEntry := range fset.Files {
		fullpath := fileEntry.FullPath()
		if strings.Contains(fullpath, common.Folder()) {
			commonFiles = append(commonFiles, fileEntry)
		}
	}

	log.Println("Found", len(commonFiles), "files")

	common.Traits.Load(commonFiles)
}
