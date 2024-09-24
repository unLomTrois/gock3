package files

import (
	"errors"
	"sync"
)

// Error to be returned if the index is out of bounds
var ErrIndexOutOfBounds = errors.New("index out of bounds")

type PathTableIndex struct {
	index uint32
}

type PathTable interface {
	Store(fullpath string) PathTableIndex
}

type PathTableStore struct {
	fullpath string
}

// Singleton of PathTable
type pathTable struct {
	paths []PathTableStore
	mu    sync.RWMutex
}

type PathTableStatic struct{}

var (
	pathTableInstance *pathTable
	once              sync.Once
	PATHTABLE         PathTableStatic
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
func (PathTableStatic) Store(fullpath string) *PathTableIndex {
	return GetPathTableInstance().store(fullpath)
}

func (pt *pathTable) store(fullpath string) *PathTableIndex {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	idx := &PathTableIndex{index: uint32(len(pt.paths))}
	pt.paths = append(pt.paths, PathTableStore{fullpath: fullpath})
	return idx
}

// Public LookupFullpath method that calls the private method after getting the singleton instance.
func (PathTableStatic) LookupFullpath(index PathTableIndex) (string, error) {
	return GetPathTableInstance().lookupFullpath(index)
}

// Private lookupFullpath method with RLock and RUnlock for thread-safe read access.
func (pt *pathTable) lookupFullpath(index PathTableIndex) (string, error) {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	if index.index >= uint32(len(pt.paths)) {
		return "", ErrIndexOutOfBounds
	}

	return pt.paths[index.index].fullpath, nil
}

// ResetPathTable is a helper function to reset the singleton for testing purposes.
func resetPathTable() {
	pathTableInstance = GetPathTableInstance()
	pathTableInstance.paths = make([]PathTableStore, 0)
}
