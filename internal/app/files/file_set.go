package files

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type FileSet struct {
	// path to ck3/game
	VanillaRoot string
	Mod         *ModLoader
	Files       []*FileEntry
}

type ModLoader struct {
	// path to ck3/mod
	Root         string
	ReplacePaths []string
}

func NewFileSet(vanillaRoot string, mod *ModLoader) *FileSet {
	return &FileSet{
		VanillaRoot: vanillaRoot,
		Mod:         mod,
	}
}

func NewModLoader(modRoot string, replacePaths []string) *ModLoader {
	return &ModLoader{
		Root:         modRoot,
		ReplacePaths: replacePaths,
	}
}

func (fset *FileSet) Scan(path string) error {
	// walk dir and add to fset.Files
	// use filepath.Walk

	clean_path := filepath.Clean(path)

	err := filepath.Walk(clean_path, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			inner_path := strings.TrimPrefix(p, clean_path)
			fmt.Println(p)
			fmt.Println(inner_path)

			file_entry := NewFileEntry(p, Mod)

			fset.Files = append(fset.Files, file_entry)
		} else {
		}

		return nil
	})

	return err
}
