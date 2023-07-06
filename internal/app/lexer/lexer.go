package lexer

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Lexer struct {
	Text    []byte
	cursor  int
	_string []byte
	Line    int
}

func NormalizeText(text []byte) []byte {

	text = bytes.TrimSpace(text)
	// text = bytes.ReplaceAll(text, []byte("    "), []byte("\t"))
	text = bytes.ReplaceAll(text, []byte("\r\n"), []byte("\n"))
	// text = bytes.ReplaceAll(text, []byte("\t\n"), []byte("\n"))
	// text = bytes.ReplaceAll(text, []byte("\t"), []byte(""))
	// text = bytes.ReplaceAll(text, []byte(" = "), []byte("="))
	// text = bytes.ReplaceAll(text, []byte("= {"), []byte("={"))

	// replace \n\n\n.. with \n\n
	reg := regexp.MustCompile(`\n{3,}`)
	text = reg.ReplaceAll(text, []byte("\n\n"))

	return text
}

func New(text []byte) *Lexer {
	normalized := NormalizeText(text)
	// fmt.Println(strconv.Quote(string(normalized)))

	new_file, _ := os.Create("./tmp/normalized.txt")
	defer new_file.Close()

	w := bufio.NewWriter(new_file)
	w.Write(normalized)
	w.Flush()

	return &Lexer{
		Text:   normalized,
		cursor: 0,
		Line:   1,
	}
}

func (l *Lexer) hasMoreTokens() bool {
	return l.cursor < len(l.Text)
}

// func (l *Lexer) isEOF() bool {
// 	return l.cursor == len(l.Text)
// }

func (l *Lexer) match(reg *regexp.Regexp, text []byte) []byte {
	if match := reg.Find(text); match != nil {
		l.cursor += len(match)
		return match
	}
	return nil
}

func (l *Lexer) GetNextToken() (*Token, error) {
	if !l.hasMoreTokens() {
		return nil, nil
	}

	l._string = l.Text[l.cursor:]

	for k, token_type := range Spec {
		// todo: implement less greedy matching
		// fmt.Println("try:", k, "on: ", string(l._string[0:10]))
		reg := regexp.MustCompile(k)
		token_value := l.match(reg, l._string)
		if token_value == nil {
			// fmt.Println("continue")
			continue
		}
		if token_type == WORD {
			match, err := regexp.MatchString(`^scripted_(trigger|effect)`, string(l._string))
			if err != nil {
				return nil, err
			}
			if match {
				token_type = SCRIPT
			}
		}

		if token_type == NEXTLINE {
			l.Line++
			return l.GetNextToken()
		}
		if token_type == NULL {
			// fmt.Println("null")
			return l.GetNextToken()
		}
		// fmt.Println(string(token_type), strconv.Quote(string(token_value)))

		// fmt.Println("return")
		return &Token{
			Type:  token_type,
			Value: token_value,
		}, nil
	}

	log.Println(string(l._string[0]))

	panic("[Lexer] Unexpected token: " + strconv.Quote(string(l._string[0])))
}

func (l *Lexer) GetContext(window int) string {
	if l.cursor < len(l.Text) {
		end := l.cursor + window
		if end > len(l.Text) {
			end = len(l.Text)
		}
		return string(l.Text[l.cursor:end])
	}
	return ""
}
