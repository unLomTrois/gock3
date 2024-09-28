package validator

import (
	"fmt"

	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

type BlockValidator struct {
	block  *ast.FieldBlock
	fields map[string]*ast.Field
	*report.ErrorManager
}

func NewBlockValidator(block *ast.FieldBlock) *BlockValidator {
	bv := &BlockValidator{
		block:        block,
		fields:       make(map[string]*ast.Field),
		ErrorManager: report.NewErrorManager(),
	}
	bv.buildFieldMap()
	return bv
}

func (bv *BlockValidator) buildFieldMap() {
	for _, field := range bv.block.Values {
		key := field.Key.Value

		bv.fields[key] = field
	}
}

func (bv *BlockValidator) ExpectBlock(key string) ast.Block {
	field, exists := bv.fields[key]
	if !exists {
		return nil
	}

	block, ok := field.Value.(ast.Block)
	if !ok {
		return nil
	}

	return block
}

func (bv *BlockValidator) ExpectToken(key string) *tokens.Token {
	field, exists := bv.fields[key]
	if !exists {
		return nil
	}

	token, ok := field.Value.(*tokens.Token)
	if !ok {
		err := report.FromToken(token, severity.Error, "expected a token, not a block")
		bv.AddError(err)

		return nil
	}

	return token
}

func (bv *BlockValidator) ExpectInteger(key string) bool {
	token := bv.ExpectToken(key)
	if token == nil {
		return false
	}

	ok := token.IsType(tokens.NUMBER)
	if !ok {
		err := report.FromToken(token, severity.Error, "expected integer")
		bv.AddError(err)
	}

	return ok
}

func (bv *BlockValidator) ExpectString(key string) bool {
	token := bv.ExpectToken(key)
	if token == nil {
		return false
	}

	ok := token.IsType(tokens.QUOTED_STRING)
	if !ok {
		err := report.FromToken(token, severity.Error, "expected string")
		bv.AddError(err)
	}

	return ok
}

func (bv *BlockValidator) RequireField(key string) bool {
	if _, exists := bv.fields[key]; !exists {
		err := report.FromBlock(bv.block, severity.Error, fmt.Sprintf("required field '%s' is missing", key))
		bv.AddError(err)
		return false
	}
	return true
}
