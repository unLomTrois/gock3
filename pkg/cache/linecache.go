package cache

import "github.com/unLomTrois/gock3/internal/app/files"

type LineCache struct {
	cache map[files.PathTableIndex][]string
}

func NewLineCache() *LineCache {
	return &LineCache{
		cache: make(map[files.PathTableIndex][]string),
	}
}

func (l *LineCache) Get(path files.PathTableIndex) ([]string, bool) {
	line, ok := l.cache[path]
	return line, ok
}

func (l *LineCache) Set(path files.PathTableIndex, line []string) {
	l.cache[path] = line
}
