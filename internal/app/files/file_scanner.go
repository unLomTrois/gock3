package files

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

// этот модуль будет сканировать все файлы в двух директориях, в директории игры и директории мода (включая вложенные файлы в другие директории)
// если возникает колизия названий файлов, приоритет будет у мода, записан будет только один файл
// также, будет аргумент замены путей, который будет давать приоритет моду над целой субдиректорией
func Scan(gameFolder string, modFolder string, replacePaths []string) ([]*FileEntry, error) {
	log.Println("Scanning game folder", gameFolder)

	cleanReplacePaths := make([]string, 0, len(replacePaths))
	for _, rp := range replacePaths {
		cleanReplacePaths = append(cleanReplacePaths, filepath.Clean(rp))
	}

	var fileMap = make(map[string]*FileEntry)

	err := filepath.WalkDir(gameFolder, func(subpath string, d fs.DirEntry, err error) error {
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

		fileName := filepath.Base(subpath)
		fileEntry := NewFileEntry(subpath, FileKind(Vanilla))
		fileMap[fileName] = fileEntry

		return nil
	})
	if err != nil {
		return nil, err
	}

	log.Println("Scanning mod folder", modFolder)
	err = filepath.WalkDir(modFolder, func(subpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !(strings.HasSuffix(subpath, ".txt")) {
			return nil
		}

		fileName := filepath.Base(subpath)
		fileEntry := NewFileEntry(subpath, FileKind(Mod))
		fileMap[fileName] = fileEntry

		return nil
	})
	if err != nil {
		return nil, err
	}

	var fileEntries []*FileEntry

	for _, fileEntry := range fileMap {
		fileEntries = append(fileEntries, fileEntry)
	}

	log.Println("Found", len(fileEntries), "files")

	return fileEntries, err
}
