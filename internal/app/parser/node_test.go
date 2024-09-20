package parser

import (
	"fmt"
	"testing"

	"github.com/unLomTrois/lexiCK3/internal/app/tokens"
)

func TestBlockField__Value(t *testing.T) {

	// we parse this string:
	// namespace = test
	// event = {
	// 	type = character_event
	// }

	b := FileBlock{
		Values: []*Field{
			{
				Key:      &tokens.Token{Value: "namespace", Type: tokens.WORD},
				Operator: &tokens.Token{Value: "=", Type: tokens.EQUALS},
				Value:    &tokens.Token{Value: "test", Type: tokens.WORD},
			},
			{
				Key:      &tokens.Token{Value: "event", Type: tokens.WORD},
				Operator: &tokens.Token{Value: "=", Type: tokens.EQUALS},
				Value: FieldBlock{
					Values: []*Field{
						{
							Key:      &tokens.Token{Value: "type", Type: tokens.WORD},
							Operator: &tokens.Token{Value: "=", Type: tokens.EQUALS},
							Value:    &tokens.Token{Value: "character_event", Type: tokens.WORD},
						},
					},
				},
			},
		},
	}

	// color = { 255 255 255 }

	color := FileBlock{
		Values: []*Field{
			{
				Key:      &tokens.Token{Value: "new", Type: tokens.WORD},
				Operator: &tokens.Token{Value: "=", Type: tokens.EQUALS},
				Value: TokenBlock{
					Values: []*tokens.Token{
						{Value: "255", Type: tokens.WORD},
						{Value: "255", Type: tokens.WORD},
						{Value: "255", Type: tokens.WORD},
					},
				},
			},
		},
	}

	fmt.Println(b)

	kek := color.Values[0].Value.(TokenBlock).Values[0]
	fmt.Println(kek)
}
