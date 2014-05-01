
package main

import (
	"io"
	"io/ioutil"
	//"fmt"
	//"strings"
	//"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	tIndent tokenType = iota
	tDedent
	tKeyword
	tDelimiter
	tOperator
	tIdentifier
	tLiteral
)

type token struct {
	typ tokenType
	value string
}

type lexer struct {
	input string
	start int
	index int
	indent []int
	tokens chan token
}

func lex(r io.Reader) (chan token, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil { return nil, err }
	l := &lexer{
		tokens: make(chan token),
	}
	l.input = string(b)
	go l.run()
	return l.tokens, nil
}

func (l *lexer) run() {
	defer close(l.tokens)
	for {
		r, s := utf8.DecodeRuneInString(l.input[l.index:])
		if r == utf8.RuneError && s == 1 { return }


		l.index += s
		if l.index == len(l.input) - 1 { return }
	}
}

