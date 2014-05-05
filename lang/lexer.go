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

type stateFn func(*lexer) stateFn

type lexer struct {
	input		string
	tokens      chan token
	indent      []int
	indent_rune rune
	start		int
	pos			int
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
		for state := lexIndent; state != nil; {
			state = state(l)
		}
	}()
	return l.tokens
}

func (l *lexer) next() (r rune) {
	r, s := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += s
	if r == '\ufffd' && s == 1 {
		log.Fatal("Invalid input encoding")
	}
	return
}

func (l *lexer) lexIndent() stateFn {

}

func (l *lexer) lexIdentifer() stateFn {

}

func (l *lexer) lexLiteral() stateFn {

}
