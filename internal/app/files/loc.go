package files

import (
	"path/filepath"
)

// Предполагаем, что PathTableIndex, FileKind, и MacroMapIndex уже определены

// Loc представляет позицию сущности в файле
type Loc struct {
	idx    PathTableIndex // Индекс в PathTable
	line   uint32         // Номер строки (0 означает весь файл)
	column uint16         // Номер столбца
	kind   FileKind       // Тип файла
}

// ForFile создает новый Loc для файла
func ForFile(pathname string, kind FileKind, fullpath string) Loc {
	idx := PATHTABLE.Store(pathname, fullpath)
	return Loc{
		idx:    idx,
		kind:   kind,
		line:   0,
		column: 0,
	}
}

// Filename возвращает имя файла из Loc
func (loc *Loc) Filename() (string, error) {
	path, err := PATHTABLE.LookupPath(loc.idx)
	if err != nil {
		return "", err
	}
	filename := filepath.Base(path)
	return filename, nil
}

// Pathname возвращает относительный путь из Loc
func (loc *Loc) Pathname() (string, error) {
	path, err := PATHTABLE.LookupPath(loc.idx)
	if err != nil {
		return "", err
	}
	return path, nil
}

// Fullpath возвращает полный путь из Loc
func (loc *Loc) Fullpath() (string, error) {
	fullpath, err := PATHTABLE.LookupFullpath(loc.idx)
	if err != nil {
		return "", err
	}
	return fullpath, nil
}

// SameFile проверяет, ссылается ли Loc на тот же файл, что и другой Loc
func (loc *Loc) SameFile(other Loc) bool {
	return loc.idx == other.idx
}

// LocFromFileEntry создает Loc из FileEntry
func LocFromFileEntry(entry *FileEntry) (Loc, error) {
	if entry.PathIdx() != nil {
		return Loc{
			idx:    *entry.PathIdx(),
			kind:   entry.Kind(),
			line:   0,
			column: 0,
		}, nil
	} else {
		err := entry.StoreInPathTable()
		if err != nil {
			return Loc{}, err
		}
		return Loc{
			idx:    *entry.PathIdx(),
			kind:   entry.Kind(),
			line:   0,
			column: 0,
		}, nil
	}
}
