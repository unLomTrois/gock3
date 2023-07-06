package parser

import (
	"ck3-parser/internal/app/lexer"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Parser struct {
	Filepath  string       `json:"filepath"`
	lexer     *lexer.Lexer `json:"-"`
	lookahead *lexer.Token `json:"-"`
	Data      []*Node      `json:"data"`
	scope     *Node
}

func New(file *os.File) (*Parser, error) {
	file_path, err := filepath.Abs(file.Name())
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	lexer := lexer.New(b)

	return &Parser{
		Filepath:  file_path,
		lexer:     lexer,
		lookahead: nil,
		Data:      nil,
		scope:     nil,
	}, nil
}

func (p *Parser) Parse() (*Parser, error) {
	var err error
	p.lookahead, err = p.lexer.GetNextToken()
	if err != nil {
		return nil, err
	}

	p.Data = p.NodeList()

	return p, nil
}

func (p *Parser) NodeList(opt_stop_lookahead ...lexer.TokenType) []*Node {
	nodes := make([]*Node, 0)

	for {
		if p.lookahead == nil {
			break
		}
		if len(opt_stop_lookahead) > 0 && p.lookahead.Type == opt_stop_lookahead[0] {
			break
		}

		new_node := p.Node()
		nodes = append(nodes, new_node)
	}

	return nodes
}

func (p *Parser) Node() *Node {
	fmt.Println("\n[LINE]", p.lexer.Line)

	if p.scope != nil {
		fmt.Println("[SCOPE]", p.scope.Key)
	}

	switch p.lookahead.Type {
	case lexer.SCRIPT:
		return p.ScriptNode()
	case lexer.COMMENT:
		return p.CommentNode()
	default:
		return p.ExpressionNode()
	}
}

func (p *Parser) CommentNode() *Node {
	fmt.Println("[COMMENT-NODE]")

	return &Node{
		Type: Comment,
		Data: p.CommentLiteral(),
	}
}

func (p *Parser) ScriptNode() *Node {
	fmt.Println("[SCRIPT-NODE]")
	p.ScriptLiteral()

	fmt.Println("--[KEY]")
	key := p.Literal()
	fmt.Println("--[OPERATION]")
	operator := p.expect(lexer.EQUALS)
	node := &Node{
		Type:     Script,
		Key:      key,
		Operator: string(operator.Value),
		Data:     nil,
	}
	if p.scope == nil {
		node.Type = Entity
		p.scope = node
		fmt.Println("--[NEW SCOPE]", p.scope.Key)
	}

	node.Data = p.BlockNode()

	fmt.Println("--[ENDSCRIPT]", p.scope.Key)

	p.expect(lexer.END)
	if p.scope == node {
		p.scope = nil
	}

	return node
}

func (p *Parser) ExpressionNode() *Node {
	fmt.Println("[EXPRESSION-NODE]")
	fmt.Println("--[KEY]")
	key := p.Literal()

	var _type NodeType
	var _operator *lexer.Token
	var _opvalue string
	fmt.Println("--[OPERATION]")
	switch p.lookahead.Type {
	case lexer.EQUALS:
		_operator = p.expect(lexer.EQUALS)
		if string(_operator.Value) == "==" {
			_type = Comparison
		} else {
			_type = Property
		}
		_opvalue = string(_operator.Value)
	case lexer.COMPARISON:
		_operator = p.expect(lexer.COMPARISON)
		_type = Comparison
		_opvalue = string(_operator.Value)
	}

	var value interface{}

	switch p.lookahead.Type {
	case lexer.WORD, lexer.NUMBER:
		fmt.Println("--[VALUE]")
		value = p.Literal()
		return &Node{
			Type:     _type,
			Key:      key,
			Operator: _opvalue,
			Data:     value,
		}
	case lexer.START:
		// fmt.Println("--[BLOCK]")
		node := &Node{
			Type:     Block,
			Key:      key,
			Operator: _opvalue,
			Data:     nil,
		}

		if p.scope == nil {
			node.Type = Entity
			p.scope = node
			fmt.Println("--[NEW SCOPE]", p.scope.Key)
		}

		node.Data = p.BlockNode()
		p.expect(lexer.END)

		if p.scope == node {
			p.scope = nil
		}

		return node
	default:
		return nil
	}
}

func (p *Parser) BlockNode() []*Node {
	fmt.Println("--[BlockNode]")

	p.expect(lexer.START)

	if p.lookahead.Type == lexer.END {
		// 	p.expect(lexer.END)

		return nil
	} else {
		return p.NodeList(lexer.END)
	}
}

func (p *Parser) Literal() interface{} {
	switch p.lookahead.Type {
	// case lexer.SCRIPT:
	// 	return p.ScriptLiteral()
	case lexer.WORD:
		return p.WordLiteral()
	case lexer.NUMBER:
		return p.NumberLiteral()
	case lexer.COMMENT:
		return p.CommentLiteral()
	default:
		log.Println("[Parser]", p.lookahead.Value, p.lookahead.Type, p.lexer.Line)
		panic(fmt.Sprintf("[Parser] Unexpected Literal: %q, with type of: %s and line: %d\ncontext: %s",
			p.lookahead.Value, p.lookahead.Type, p.lexer.Line, p.lexer.GetContext(100)))
	}
}

func (p *Parser) ScriptLiteral() string {
	token := p.expect(lexer.SCRIPT)
	return string(token.Value)
}

func (p *Parser) WordLiteral() string {
	token := p.expect(lexer.WORD)
	return string(token.Value)
}

func (p *Parser) NumberLiteral() float32 {
	token := p.expect(lexer.NUMBER)
	number, err := strconv.ParseFloat(string(token.Value), 32)
	if err != nil {
		panic("[Parser] Unexpected NumberLiteral: " + strconv.Quote(string(token.Value)))
	}
	return float32(number)
}

func (p *Parser) CommentLiteral() string {
	token := p.expect(lexer.COMMENT)
	return string(token.Value)
}

func (p *Parser) expect(tokentype lexer.TokenType) *lexer.Token {
	token := p.lookahead
	if token == nil {
		panic("[Parser] Unexpected end of input, expected: " + string(tokentype))
	}
	if token.Type != tokentype {
		panic("[Parser] Unexpected token: \"" + string(token.Value) + "\" with type of " + string(token.Type) + ", expected type: " + string(tokentype))
	}
	fmt.Println(string(p.lookahead.Type), strconv.Quote(string(p.lookahead.Value)), p.lexer.Line)

	var err error
	p.lookahead, err = p.lexer.GetNextToken()
	if err != nil {
		panic(err)
	}
	return token
}
