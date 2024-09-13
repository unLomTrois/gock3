package files

import (
	"sync"
)

type PathTableIndex struct {
	index uint32
}

type PathTable interface {
	Store(local string, fullpath string) PathTableIndex
}

type PathTableStore struct {
	local    string
	fullpath string
}

// Singleton of PathTable
type pathTable struct {
	paths []PathTableStore
	mu    sync.RWMutex
}

var (
	pathTableInstance *pathTable
	once              sync.Once
	PATHTABLE         pathTable
)

func GetPathTableInstance() *pathTable {
	once.Do(func() {
		pathTableInstance = &pathTable{
			paths: make([]PathTableStore, 0),
			mu:    sync.RWMutex{},
		}
	})
	return pathTableInstance
}

// Store is usually called from PATHTABLE, which is a "static namespace" of the package.
// Hence public Store method should call GetPathTableInstance to get an actual singleton, and call its private store method
func (pt *pathTable) Store(local string, fullpath string) PathTableIndex {
	return GetPathTableInstance().store(local, fullpath)
}

func (pt *pathTable) store(local string, fullpath string) PathTableIndex {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	idx := PathTableIndex{index: uint32(len(pt.paths))}
	pt.paths = append(pt.paths, PathTableStore{local: local, fullpath: fullpath})
	return idx
}

// ResetPathTable is a helper function to reset the singleton for testing purposes.
func resetPathTable() {
	pathTableInstance = GetPathTableInstance()
	pathTableInstance.paths = make([]PathTableStore, 0)
}
