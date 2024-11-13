package files

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type FileSet struct {
	// path to ck3/game
	VanillaRoot string
	ModLoader   *ModLoader
	Files       []*FileEntry
}

type ModLoader struct {
	// path to ck3/mod
	Root         string
	ReplacePaths []string
}

func NewFileSet(vanillaRoot string, modLoader *ModLoader) *FileSet {
	return &FileSet{
		VanillaRoot: vanillaRoot,
		ModLoader:   modLoader,
	}
}

func NewModLoader(modRoot string, replacePaths []string) *ModLoader {
	return &ModLoader{
		Root:         modRoot,
		ReplacePaths: replacePaths,
	}
}

func (fset *FileSet) Scan(path string) error {
	log.Println("Scanning", path)

	cleanReplacePaths := make([]string, 0, len(fset.ModLoader.ReplacePaths))
	for _, replacePath := range fset.ModLoader.ReplacePaths {
		cleanReplacePaths = append(cleanReplacePaths, filepath.Clean(replacePath))
	}

	err := filepath.WalkDir(path, func(subpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			for _, replacePath := range cleanReplacePaths {
				if strings.Contains(subpath, replacePath) {
					return filepath.SkipDir
				}
			}

			return nil
		}

		if !(strings.HasSuffix(subpath, ".txt")) {
			return nil
		}

		file_entry := NewFileEntry(subpath, Mod)

		fset.Files = append(fset.Files, file_entry)

		return nil
	})

	log.Println("Found", len(fset.Files), "files")

	return err
}
