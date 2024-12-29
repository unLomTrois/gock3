package files

type FileSet struct {
	// path to ck3/game
	VanillaRoot string
	ModLoader   *ModLoader
	Files       []*FileEntry
}

type ModLoader struct {
	// path to ck3/mod
	Root         string
	ReplacePaths []string
}

func NewFileSet(vanillaRoot string, modLoader *ModLoader) *FileSet {
	return &FileSet{
		VanillaRoot: vanillaRoot,
		ModLoader:   modLoader,
	}
}

func NewModLoader(modRoot string, replacePaths []string) *ModLoader {
	return &ModLoader{
		Root:         modRoot,
		ReplacePaths: replacePaths,
	}
}
