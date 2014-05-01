
package main

import (
	"io"
	"bufio"
	"log"
	//"fmt"
	//"strings"
	"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	tNewline tokenType = iota
	tIndent
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
	line string
	lineNum int
	start int
	index int
	next int
	indent []int
	indent_char rune
	tokens chan token
}

func lex(r io.Reader) (chan token) {
	s := bufio.NewScanner(r)
	l := &lexer{
		tokens: make(chan token),
		indent:	[]int{ 0 },
	}
	go func() {
		defer close(l.tokens)
		for s.Scan() {
			l.lineNum++
			l.line = s.Text()
			l.start = 0
			l.index = 0
			l.next = 0
			l.lexLine()
			l.tokens <- token{ typ: tNewline }
		}
		if err := s.Err(); err != nil {
			log.Fatalf("Scanning error at line #%d:\n\t%s", l.lineNum, err)
		}
	}()
	return l.tokens
}

func (l *lexer) lexLine() {
	//var t token
	for i := 0; i < len(l.line); i++ {
		r, s := utf8.DecodeRuneInString(l.line[l.index:])
		if r == utf8.RuneError && s == 1 {
			log.Fatal("Invalid input encoding")
		}
		l.next += s
		if r == '#' { return }
		if l.start == 0 {
			l.parseIndent(r)
		}
		l.index = l.next
	}
}

func (l *lexer) parseIndent(r rune) {
	if !unicode.IsSpace(r) {
		cur := l.indent[len(l.indent) - 1]
		if l.index > cur {
			l.tokens <- token{ typ: tIndent }
			l.indent = append(l.indent, l.index)
		} else if l.index < cur {
			for l.indent[len(l.indent) - 1] != l.index {
				l.tokens <- token{ typ: tDedent }
				l.indent = l.indent[:len(l.indent) - 1]
				if l.index > l.indent[len(l.indent) - 1] {
					log.Println("Line #%d: Inconsistent indenting", l.lineNum)
					l.indent[len(l.indent) - 1] = l.index
				}
			}
		}
		l.start = l.next
		return
	}
	if l.indent_char == 0 {
		if r != '\t' && r != ' ' {
			log.Fatalf("Line #%d: Invalid indentation character (must be ' ' or '\\t')", l.lineNum)
		}
		l.indent_char = r
	} else if r != l.indent_char {
		log.Fatal("Line #%d: Mixing tabs and spaces for indentation", l.lineNum)
	}
}
