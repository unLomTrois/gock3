package files

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

// Scan scans two directories (the game folder and the mod folder) to find .txt files,
// including those in subdirectories. If there is a filename collision, the file in the mod folder
// takes precedence. A list of paths (replacePaths) can be provided so that any matching
// subdirectory under the game folder will be skipped, effectively giving priority to the mod folder
// for that subdirectory.
func Scan(gameFolder string, modFolder string, replacePaths []string) ([]*FileEntry, error) {
	// Normalize replacePaths to clean directory paths
	normalizedReplacePaths := make([]string, 0, len(replacePaths))
	for _, path := range replacePaths {
		normalizedReplacePaths = append(normalizedReplacePaths, filepath.Clean(path))
	}

	// A map to hold files uniquely by their filename
	fileMap := make(map[string]*FileEntry)

	// Helper function to handle skipping directories and adding files
	scanFunc := func(root string, kind FileKind) fs.WalkDirFunc {
		return func(subpath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Skip entire directories if they match any replacePath
			if d.IsDir() && kind == Vanilla {
				for _, replacePath := range normalizedReplacePaths {
					if strings.Contains(subpath, replacePath) {
						return filepath.SkipDir
					}
				}
				return nil
			}

			// Only consider .txt files
			if !strings.HasSuffix(subpath, ".txt") {
				return nil
			}

			filename := filepath.Base(subpath)
			fileEntry := NewFileEntry(subpath, kind)
			fileMap[filename] = fileEntry

			return nil
		}
	}

	log.Printf("Scanning game folder: %s\n", gameFolder)
	err := filepath.WalkDir(gameFolder, scanFunc(gameFolder, Vanilla))
	if err != nil {
		return nil, err
	}

	log.Printf("Scanning mod folder: %s\n", modFolder)
	if err := filepath.WalkDir(modFolder, scanFunc(modFolder, Mod)); err != nil {
		return nil, err
	}

	// Collect file entries from the map
	var fileEntries []*FileEntry
	for _, entry := range fileMap {
		fileEntries = append(fileEntries, entry)
	}

	log.Printf("Found %d files\n", len(fileEntries))
	return fileEntries, nil
}
