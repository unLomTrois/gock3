package ast

import "github.com/unLomTrois/gock3/internal/app/lexer/tokens"

type Field struct {
	Key      *tokens.Token `json:"key"`
	Operator *tokens.Token `json:"operator"`
	Value    BV            `json:"value"`
}
