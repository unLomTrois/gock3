package linter

import (
	"ck3-parser/internal/app/parser"
	"os"
)

type IntendStyle string

const (
	TABS   IntendStyle = "TABS"
	SPACES IntendStyle = "SPACES"
)

type Linter struct {
	ParseTree   []*parser.Node `json:"tree"`
	Level       int            `json:"level"`
	towrite     []byte
	intendstyle IntendStyle
}

func New(parsetree []*parser.Node) *Linter {
	return &Linter{
		ParseTree:   parsetree,
		Level:       0,
		towrite:     []byte{},
		intendstyle: "TAB",
	}
}

func (l *Linter) Lint() []byte {
	for _, node := range l.ParseTree {
		l.lintNode(node)
		// if i != len(l.ParseTree)-1 {
		// 	l.towrite = append(l.towrite, byte('\n'))
		// }
	}

	return l.towrite
}

func (l *Linter) lintNode(node *parser.Node) {
	switch node.Type {
	case parser.Comment:
		l.lintComment(node)
	case parser.Property:
		l.lintProperty(node)
	case parser.Block:
		l.lintBlock(node)
	}
}

func (l *Linter) lintComment(node *parser.Node) {
	if len(l.towrite) > 0 && l.towrite[len(l.towrite)-1] != ' ' {
		l.intend()
	}
	l.towrite = append(l.towrite, []byte(node.Value.(string))...)
	l.nextLine()
	// if l.singleline {
	// 	l.towrite = append(l.towrite, byte(' '))
	// } else if l.Level != 0 {
	// 	l.towrite = append(l.towrite, byte('\n'))
	// }
}

func (l *Linter) lintProperty(node *parser.Node) {
	if len(l.towrite) > 0 && l.towrite[len(l.towrite)-1] != ' ' {
		l.intend()
	}

	l.towrite = append(l.towrite, node.KeyLiteral()...)
	l.equals()
	l.towrite = append(l.towrite, node.DataLiteral()...)

	l.nextLine()
}

func (l *Linter) lintBlock(node *parser.Node) {
	if len(l.towrite) > 0 && l.towrite[len(l.towrite)-1] != ' ' {
		l.intend()
	}

	children := node.Value.([]*parser.Node)
	l.Level++

	l.towrite = append(l.towrite, node.KeyLiteral()...)
	l.equals()
	l.towrite = append(l.towrite, byte('{'))

	l.nextLine()

	for _, c := range children {
		l.lintNode(c)
	}
	l.Level--
	l.intend()
	l.towrite = append(l.towrite, byte('}'))

	l.nextLine()
}

func (l *Linter) intend() {
	i := 0
	for i < l.Level {
		l.towrite = append(l.towrite, []byte("  ")...)
		i++
	}
}

func (l *Linter) equals() {
	l.towrite = append(l.towrite, byte(' '))
	l.towrite = append(l.towrite, byte('='))
	l.towrite = append(l.towrite, byte(' '))
}

func (l *Linter) nextLine() {
	l.towrite = append(l.towrite, byte('\n'))
}

func (l *Linter) Save(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(l.towrite)
	if err != nil {
		return err
	}

	return nil
}
