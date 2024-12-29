package data

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/unLomTrois/gock3/internal/app/files"
)

type History struct {
	Characters *HistoryCharacters
}

func NewHistory() *History {
	return &History{
		Characters: NewHistoryCharacters(),
	}
}

// Folder returns the folder path for common, using the correct
// path separator for the operating system.
func (c *History) Folder() string {
	return filepath.Join("game", "history")
}

func (history *History) Load(fset *files.FileSet) []Entity {
	var files []*files.FileEntry

	for _, fileEntry := range fset.Files {
		fullpath := fileEntry.FullPath()
		if strings.Contains(fullpath, history.Folder()) {
			files = append(files, fileEntry)
		}

		if strings.Contains(fullpath, filepath.Clean(fset.ModLoader.Root)) {
			files = append(files, fileEntry)
		}
	}

	var entities []Entity

	log.Printf("Found %d history files", len(files))
	characters := history.Characters.Load(files)

	for _, character := range characters {
		entities = append(entities, character)
	}

	return entities
}
