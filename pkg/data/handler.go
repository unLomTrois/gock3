package data

import "github.com/unLomTrois/gock3/internal/app/parser/ast"

type DataHandler interface {
	DataFolder
	DataLoader
}

type DataFolder interface {
	Folder() string
}

type DataLoader interface {
	Load() *ast.AST
}
