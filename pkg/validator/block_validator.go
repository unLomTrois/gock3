package validator

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
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

// NewBlockValidator creates a new instance of BlockValidator for the given FieldBlock.
// It initializes the fields map by converting the block's values into a map
// for efficient key-based lookups and sets up an ErrorManager for reporting validation errors.
//
// Parameters:
//   - block: A pointer to the FieldBlock containing fields to validate.
//
// Returns:
//   - A pointer to the newly created BlockValidator instance.
func NewBlockValidator(block *ast.FieldBlock) *BlockValidator {
	fields := make(map[string]*ast.Field, len(block.Values))
	for _, field := range block.Values {
		fields[field.Key.Value] = field
	}

	return &BlockValidator{
		block:        block,
		fields:       fields,
		ErrorManager: report.NewErrorManager(),
	}
}

// ExpectBlock checks that there is a field with a certain key whose value is a block.
func (bv *BlockValidator) ExpectBlock(key string) ast.Block {
	field, exists := bv.fields[key]
	if !exists {
		return nil
	}

	block, ok := field.Value.(ast.Block)
	if !ok {
		err := report.FromToken(field.Key, severity.Error, "expected a block, not a token")
		bv.AddError(err)
		return nil
	}

	return block
}

// ExpectToken checks that there is a field with a certain key whose value is a token.
func (bv *BlockValidator) ExpectToken(key string) *tokens.Token {
	field, exists := bv.fields[key]
	if !exists {
		return nil
	}

	token, ok := field.Value.(*tokens.Token)
	if !ok {
		err := report.FromLoc(field.Key.Loc, severity.Error, "expected a token, not a block")
		bv.AddError(err)
		return nil
	}

	return token
}

// expect value to be, i.e. token.Value ==
func (bv *BlockValidator) ExpectValueToBe(key string, value string) bool {
	token := bv.ExpectToken(key)
	if token == nil {
		return false
	}
	return token.Is(value)
}

// expect value to be, i.e. token.Value ==
func (bv *BlockValidator) ExpectValueToBeInSet(key string, set mapset.Set[string]) bool {
	token := bv.ExpectToken(key)
	if token == nil {
		return false
	}

	ok := set.Contains(token.Value)
	if !ok {
		err := report.FromToken(token, severity.Error, fmt.Sprintf("expected one of %v", set))
		bv.AddError(err)
	}

	return ok
}

// ExpectNumber checks that there is a field with a certain key whose value is a token and which type is a number
func (bv *BlockValidator) ExpectType(key string, tt tokens.TokenType) (*tokens.Token, bool) {
	token := bv.ExpectToken(key)
	if token == nil {
		return nil, false
	}

	ok := token.IsType(tt)
	if !ok {
		err := report.FromToken(token, severity.Error, fmt.Sprintf("expected type %s", tt))
		bv.AddError(err)
	}

	return token, ok
}

// is bool
func (bv *BlockValidator) ExpectBool(key string) (*tokens.Token, bool) {
	token, ok := bv.ExpectType(key, tokens.BOOL)

	return token, ok
}

// ExpectNumber checks that there is a field with a certain key whose value is a token and which type is a number
func (bv *BlockValidator) ExpectNumber(key string) bool {
	token := bv.ExpectToken(key)
	if token == nil {
		return false
	}

	ok := token.IsType(tokens.NUMBER)
	if !ok {
		err := report.FromToken(token, severity.Error, "expected a number")
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
	_, exists := bv.fields[key]
	if !exists {
		err := report.FromBlock(bv.block, severity.Error, fmt.Sprintf("required field '%s' is missing", key))
		bv.AddError(err)
	}
	return exists
}

// number in range (float)
func (bv *BlockValidator) ExpectNumberInRange(key string, min float64, max float64) bool {
	token, ok := bv.ExpectType(key, tokens.NUMBER)
	if !ok {
		return false
	}

	value, err := token.FloatValue()
	if err != nil {
		err := report.FromToken(token, severity.Error, "expected a number")
		bv.AddError(err)
		return false
	}
	if value < min || value > max {
		err := report.FromToken(token, severity.Error, fmt.Sprintf("expected number in range [%f, %f]", min, max))
		bv.AddError(err)
		return false
	}

	return true
}

// ban field
func (bv *BlockValidator) BanField(key string, because string) {
	_, exists := bv.fields[key]
	if exists {
		err := report.FromBlock(bv.block, severity.Error, fmt.Sprintf("field '%s' is not allowed, because %s", key, because))
		bv.AddError(err)
	}
}

// Deprecated: Add more suitable way.
func (bv *BlockValidator) Fields() map[string]*ast.Field {
	return bv.fields
}
