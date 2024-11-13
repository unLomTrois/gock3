package cache

import (
	"log"
	"os"
	"strings"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
)

type FileCache struct {
	cache map[files.PathTableIndex]string

	linecache *LineCache
}

func NewFileCache() *FileCache {
	return &FileCache{
		cache:     make(map[files.PathTableIndex]string),
		linecache: NewLineCache(),
	}
}

func (f *FileCache) Get(pathTableIndex files.PathTableIndex) (string, bool) {
	content, ok := f.cache[pathTableIndex]
	return content, ok
}

func (f *FileCache) Add(index files.PathTableIndex) {
	fullpath, err := files.PATHTABLE.LookupFullpath(index)
	if err != nil {
		panic(err)
	}

	// read file!
	content, err := os.ReadFile(fullpath)
	if err != nil {
		panic(err)
	}

	f.cache[index] = string(content)
}

func (f *FileCache) Set(index files.PathTableIndex, value string) {
	f.cache[index] = value
}

func (f *FileCache) GetLine(loc *tokens.Loc) string {
	index := loc.GetIdx()

	// check lines cache
	if lines, ok := f.linecache.Get(index); ok {
		return lines[loc.Line-1]
	}

	// check filecache and fill linecache
	if content, ok := f.Get(index); ok {
		// split content by \n
		lines := strings.Split(content, "\n")

		f.linecache.Set(index, lines)
		return lines[loc.Line-1]
	}

	// if nothing found, fill filecahce
	f.Add(index)
	log.Println("skip")

	// recursive call
	return f.GetLine(loc)
}
