package files

import (
	"errors"
	"os"
	"path/filepath"
)

type FileKind uint8

const (
	Vanilla FileKind = iota
	Mod
)

type FileEntry struct {
	// The full filesystem path of this entry
	fullpath string
	// Index into the PathTable (optional, using *PathTableIndex to allow nil)
	idx *PathTableIndex
	// Whether it's a vanilla or mod file
	kind FileKind
}

// NewFileEntry is the constructor for FileEntry.
// Ensures the path is valid and not empty.
func NewFileEntry(fullpath string, kind FileKind) *FileEntry {
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		panic("Invalid path: path does not exist")
	}

	return &FileEntry{
		fullpath: fullpath,
		kind:     kind,
		idx:      nil,
	}
}

// Kind returns the file kind (vanilla or mod).
func (fe *FileEntry) Kind() FileKind {
	return fe.kind
}

// FullPath returns the full filesystem path.
func (fe *FileEntry) FullPath() string {
	return fe.fullpath
}

// FileName returns the file name, ensuring it's not empty.
func (fe *FileEntry) FileName() string {
	return filepath.Base(fe.fullpath)
}

func (fe *FileEntry) StoreInPathTable() error {
	if fe.idx != nil {
		return errors.New("PathTableIndex is already set")
	}
	idx := PATHTABLE.Store(fe.fullpath)
	fe.idx = &idx
	return nil
}

// PathIdx returns the index into the PathTable if it exists, otherwise nil.
func (fe *FileEntry) PathIdx() *PathTableIndex {
	return fe.idx
}
