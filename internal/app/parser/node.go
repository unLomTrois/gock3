package parser

import "github.com/unLomTrois/lexiCK3/internal/app/tokens"

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
	IsBlock()
	IsBV()
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
