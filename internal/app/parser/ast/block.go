package ast

import (
	"fmt"

	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
)

// type that can be a Block or a Token
type BV interface {
	IsBV()
}

type Block interface {
	BV
	IsBlock()
}

// File Block is a top-level block with a list of fields
type FileBlock = FieldBlock

// Field Block is a block with a list of fields
type FieldBlock struct {
	Values []*Field `json:"fields"`
	Loc    tokens.Loc
}

func (fb *FieldBlock) IsBlock() {}
func (fb *FieldBlock) IsBV()    {}
func (fb *FieldBlock) GetValues() []*Field {
	return fb.Values
}

// GetField
func (fb *FieldBlock) GetField(key string) *Field {
	for _, field := range fb.Values {
		if field.Key.Value == key {
			return field
		}
	}

	return nil
}

func (fb *FieldBlock) GetFieldValue(key string) *tokens.Token {
	field := fb.GetField(key)
	if field == nil {
		return nil
	}

	return field.Value.(*tokens.Token)
}

// GetFields searches all fields with a certain key
func (fb *FieldBlock) GetFields(key string) []*Field {
	res := make([]*Field, 0)

	// search for fields with the given key
	for _, field := range fb.Values {
		if field.Key.Value == key {
			res = append(res, field)
		}
	}

	return res
}

func (fb *FieldBlock) GetFieldsValues(key string) []*tokens.Token {
	fields := fb.GetFields(key)
	res := make([]*tokens.Token, len(fields))

	for i, field := range fields {
		fmt.Println(field.Value)
		res[i] = field.Value.(*tokens.Token)
	}

	return res
}

func (fb *FieldBlock) GetFieldList(key string) []*tokens.Token {
	field := fb.GetField(key)
	if field == nil {
		return nil
	}

	switch field.Value.(type) {
	case *TokenBlock:
		return field.Value.(*TokenBlock).Values
	}

	return nil
}

func (fb *FieldBlock) GetFieldBlock(key string) *FieldBlock {
	// get the field
	field := fb.GetField(key)
	if field == nil {
		return nil
	}

	block, ok := field.Value.(*FieldBlock)
	if !ok {
		return nil
	}

	return block
}

func (fb *FieldBlock) GetTokenBlock(key string) *TokenBlock {
	// get the field
	field := fb.GetField(key)
	if field == nil {
		return nil
	}

	block, ok := field.Value.(*TokenBlock)
	if !ok {
		return nil
	}

	return block
}

// Token Block is a block with a list of tokens
type TokenBlock struct {
	Values []*tokens.Token `json:"tokens"`
}

func (tb *TokenBlock) IsBlock() {}
func (tb *TokenBlock) IsBV()    {}
