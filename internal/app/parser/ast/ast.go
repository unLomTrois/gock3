package ast

import "github.com/unLomTrois/gock3/internal/app/lexer/tokens"

type AST struct {
	Filename string   `json:"filename"`
	Fullpath string   `json:"fullpath"`
	Data     []*Field `json:"data"`
}

type Field struct {
	Key      *tokens.Token `json:"key"`
	Operator *tokens.Token `json:"operator"`
	Value    BV            `json:"value"`
}

// type that can be a Block or a Token
type BV interface {
	IsBV()
}

type Block interface {
	BV
	IsBlock()
}

// File Block is a top-level block with a list of fields
type FileBlock struct {
	Values []*Field `json:"properties"`
}

func (fb FileBlock) IsBlock() {}

// Field Block is a block with a list of fields
type FieldBlock struct {
	Values []*Field `json:"properties"`
}

func (fb FieldBlock) IsBlock() {}
func (fb FieldBlock) IsBV()    {}

// Token Block is a block with a list of tokens
type TokenBlock struct {
	Values []*tokens.Token `json:"tokens"`
}

func (tb TokenBlock) IsBlock() {}
func (tb TokenBlock) IsBV()    {}
