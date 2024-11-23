package data

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/unLomTrois/gock3/internal/app/lexer/tokens"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/validator"
)

type Trait struct {
	name  string
	key   *tokens.Token
	block *ast.FieldBlock
}

func NewTraitFromAST(key *tokens.Token, block *ast.FieldBlock) *Trait {
	return &Trait{
		name:  key.Value,
		key:   key,
		block: block,
	}
}

var categorySet = mapset.NewSet("personality", "education", "childhood", "commander", "winter_commander", "lifestyle", "court_type", "fame", "health")

func (trait *Trait) Validate() []*report.DiagnosticItem {
	fields := validator.NewBlockValidator(trait.block)

	// for _, field := range trait.block.Values {
	// 	ok := availableKeys.Contains(field.Key.Value)
	// 	if !ok {
	// 		err := report.FromToken(field.Key, severity.Error, fmt.Sprintf("Unknown key %q", field.Key.Value))
	// 		fields.AddError(err)
	// 	}
	// }

	fields.ExpectBool("genetic")
	if genetic := fields.ExpectValueToBe("genetic", "yes"); genetic {
		fields.BanField("random_creation_weight", "it is not allowed for genetic traits")
		fields.ExpectNumberInRange("birth", 0, 1)
		fields.ExpectNumberInRange("random_creation", 0, 1)
	} else {
		fields.ExpectNumberInRange("random_creation_weight", 0, 1)
		fields.BanField("birth", "it is not allowed for non genetic traits")
		fields.BanField("random_creation", "it is not allowed for non genetic traits")
	}

	fields.ExpectValueToBeInSet("category", categorySet)

	fields.ExpectNumber("minimum_age")
	fields.ExpectNumber("maximum_age")

	fields.ExpectNumber("stewardship")
	fields.ExpectNumber("diplomacy")
	fields.ExpectNumber("martial")
	fields.ExpectNumber("intrigue")
	fields.ExpectNumber("learning")

	fields.ExpectBool("physical")
	fields.ExpectBool("good")
	fields.ExpectBool("immortal")
	fields.ExpectBool("can_have_children")
	fields.ExpectBool("enables_inbred")

	return fields.Errors()
}

// todo: idea, we can make a map of keys and functions that shall validate this key to improve O from O(n) to O(1)
