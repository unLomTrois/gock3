package data

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/pkg/entity"
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

func (common *Common) Load(fset *files.FileSet) []entity.Entity {
	var files []*files.FileEntry

	for _, fileEntry := range fset.Files {
		fullpath := fileEntry.FullPath()
		if strings.Contains(fullpath, common.Folder()) {
			files = append(files, fileEntry)
		}

		if strings.Contains(fullpath, filepath.Clean(fset.ModLoader.Root)) {
			files = append(files, fileEntry)
		}
	}

	log.Printf("Found %d files in common", len(files))

	var entities []entity.Entity

	traits := common.Traits.Load(files)

	for _, trait := range traits {
		// fmt.Println(trait.Name(), trait.Location())
		entities = append(entities, trait)
	}

	return entities
}
