package main

import (
	"bufio"
	"io"
	"log"
	"unicode"
	"unicode/utf8"
)

type char struct {
	line int
	col  int
	rune
}

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
	typ   tokenType
	value string
}

type lexer struct {
	*log.Logger
	*bufio.Reader
	last        char
	cur         string
	indent      []int
	indent_rune rune
	tokens      chan token
}

func lex(r io.Reader) chan token {
	l := &lexer{
		log.New(os.Stderr, log.Prefix(), log.Flags()),
		bufio.NewReader(r),
		last:   char{rune: '\n'},
		tokens: make(chan token),
		indent: []int{0},
	}
	go func() {
		defer close(l.tokens)
		for {
		}
	}()
	return l.tokens
}

func (l *lexer) logError(err error, c char) {
	l.Panicf("Line #%d, Column #%d:\n\t%s\n", c.line, c.col, err)
}

func (l *lexer) next() (c char) {
	r, s, err := l.ReadRune()
	if err != nil {
		l.logError(err, l.last)
	}
	c := char{
		line: l.last.line,
		col:  l.last.col + len(string(l.last.rune)),
		rune: r,
	}
	if l.last.rune == '\n' {
		c.line++
		c.col = 1
	}
	if r == '\ufffd' && s == 1 {
		l.logError("Invalid input encoding", c)
	}
}

func (l *lexer) peek() (c char) {
	c = l.next()
	err := l.UnreadRune()
	if err != nil {
		l.logError(err, c)
	}
}

func (l *lexer) emit(typ tokenType, v string) {
	l.tokens <- token{typ: typ, value: v}
}

const keywords = []string{
	"const", "byte", "block", "func ", "...", "->", "jump", "return", "if",
}
const operators_delimiters = []string{
	":", ":=", "=", "+", "-", "*", "/", "**", ",", "(", ")",
}
