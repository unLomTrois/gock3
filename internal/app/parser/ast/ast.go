package ast

type AST struct {
	Filename string     `json:"filename"`
	Fullpath string     `json:"fullpath"`
	Block    *FileBlock `json:"data"`
}
