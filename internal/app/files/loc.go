package files

import (
	"path/filepath"
)

// Предполагаем, что PathTableIndex, FileKind, и MacroMapIndex уже определены

// Loc представляет позицию сущности в файле
type Loc struct {
	idx    PathTableIndex `json:"-"`
	Line   uint32         `json:"line"`
	Column uint16         `json:"column"`
	kind   FileKind       `json:"-"`
}

// ForFile создает новый Loc для файла
func ForFile(pathname string, kind FileKind, fullpath string) Loc {
	idx := PATHTABLE.Store(fullpath)
	return Loc{
		idx:    idx,
		kind:   kind,
		Line:   0,
		Column: 0,
	}
}

// Filename возвращает имя файла из Loc
func (loc *Loc) Filename() (string, error) {
	path, err := PATHTABLE.LookupFullpath(loc.idx)
	if err != nil {
		return "", err
	}
	filename := filepath.Base(path)
	return filename, nil
}

// Pathname возвращает относительный путь из Loc
func (loc *Loc) Pathname() (string, error) {
	path, err := PATHTABLE.LookupFullpath(loc.idx)
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
			Line:   1,
			Column: 1,
		}, nil
	} else {
		err := entry.StoreInPathTable()
		if err != nil {
			return Loc{}, err
		}
		return Loc{
			idx:    *entry.PathIdx(),
			kind:   entry.Kind(),
			Line:   1,
			Column: 1,
		}, nil
	}
}
