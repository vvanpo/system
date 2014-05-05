package main

import (
	"io"
	"io/ioutil"
	"strings"
	"log"
	"unicode"
	"unicode/utf8"
)

type terminal int

const (
	tIndent terminal = iota
	tDedent
	tIdentifier
	tLiteral
	tReserved
	tOperator
	tDelimiter
)

var reserved = [...]string{
	"const", "byte", "block", "func", "...", "->", "jump", "return", "if",
}
var operators = [...]string{
	":", ":=", "=", "+", "-", "*", "/", "**", "(", ")",
}
var delimiters = [...]string{
	",",
}

type token struct {
	terminal
	attr string
	num int64
}

type stateFn func(*lexer, rune) stateFn

type lexer struct {
	input		string
	tokens      chan token
	indent      []int
	indent_rune rune
	start		int
	pos			int
	width		int
}

func lex(r io.Reader) chan token {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal("Input error")
	}
	l := &lexer{
		input:	string(input),
		tokens: make(chan token),
		indent: []int{0},
	}
	go func() {
		defer close(l.tokens)
		for state := lexIndent; state != nil && l.pos < len(input); {
			n := l.next()
			if n == '#' {
				p := strings.IndexRune(l.input[l.pos:], '\n')
				if p == -1 { return }
				n = '\n'
				l.pos = p
				l.start = p
			}
			state = state(l, n)
		}
	}()
	return l.tokens
}

func (l *lexer) next() (r rune) {
	r, s := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = s
	l.pos += s
	if r == '\ufffd' && s == 1 {
		log.Fatal("Invalid input encoding")
	}
	return
}

// Call lexIndent after a '\n'
func lexIndent(l *lexer, r rune) stateFn {
	if !unicode.IsSpace(r) {
		indent := l.start - l.pos - l.width
		for i := range(l.indent) {
			if l.indent[i] == indent {
				if diff := len(l.indent) - 1 - i; diff != 0 {
					for j := 0; j < diff; j++ {
						l.tokens <- token{ terminal: tDedent }
					}
				}
			} else if l.indent[i] > indent {
				log.Fatal("Indentation mismatch")
			} else {
				l.indent = append(l.indent, indent)
				l.tokens <- token{ terminal: tIndent }
			}
			l.start = l.pos - l.width
			return lexIdentifier(l, r)
		}
	}
	if r == '\n' {
		l.start = l.pos
	} else if r != '\t' && r != ' ' {
		log.Fatal("Invalid indentation character")
	} else if l.indent_rune == 0 {
		l.indent_rune = r
	} else if r != l.indent_rune {
		log.Fatal("Mixing tabs and spaces for indentation")
	}
	return lexIndent
}

func lexIdentifier(l *lexer, r rune) stateFn {
	if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
		if l.pos == l.start + l.width {
			log.Fatalf("Invalid identifer character: %v", r)
		}
		if unicode.IsSpace(r) || r == ':' || r == '=' || r == ',' || r == '(' || r == ')' {
			val := l.input[l.start:(l.pos - l.width)]
			for _, k := range(reserved) {
				if val == k {
					l.tokens <- token{ terminal: tReserved, attr: val }
					l.start = l.pos - l.width
					return lexSpace(l, r)
				}
			}
			l.tokens <- token{ terminal: tIdentifier, attr: val }
			l.start = l.pos - l.width
			return lexSpace(l, r)
		}
	}
	if l.pos - l.width == l.start && unicode.IsDigit(r) {
		return lexLiteral(l, r)
	}
	return lexIdentifier
}

func lexLiteral(l *lexer, r rune) stateFn {
	return nil
}

func lexSpace(l *lexer, r rune) stateFn {
	l.start = l.pos - l.width
	if unicode.IsSpace(r) {
		if r == '\n' {
			return lexIndent
		}
		return lexSpace
	}
	if unicode.IsDigit(r) {
		return lexLiteral(l, r)
	}
	if unicode.IsLetter(r) || r == '_' {
		return lexIdentifier(l, r)
	}
	return lexSpace
}
