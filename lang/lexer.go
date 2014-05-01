package main

import (
	"bufio"
	"io"
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
	typ   tokenType
	value string
}

type lexer struct {
	line        string
	lineNum     int
	start       int
	index       int
	indent      []int
	indent_char rune
	tokens      chan token
}

func lex(r io.Reader) chan token {
	s := bufio.NewScanner(r)
	l := &lexer{
		tokens: make(chan token),
		indent: []int{0},
	}
	go func() {
		defer close(l.tokens)
		for s.Scan() {
			l.lineNum++
			l.line = s.Text()
			l.start = 0
			l.index = 0
			l.lexLine()
		}
		if err := s.Err(); err != nil {
			log.Fatalf("Scanning error at line #%d:\n\t%s", l.lineNum, err)
		}
	}()
	return l.tokens
}

func (l *lexer) peek(i int) (r rune, s int) {
	r, s := utf8.DecodeRuneInString(l.line[i:])
	if r == utf8.RuneError && s == 1 {
		log.Fatal("Invalid input encoding")
	}
}

func (l *lexer) lexLine() {
	var t token
	for i := 0; i < len(l.line); i++ {
		r, s := l.peek(l.index)
		if r == '#' {
			return
		}
		if l.start == 0 {
			l.parseIndent(r)
		} else {
			t = l.tokenize(r, t)
		}
		l.index += s
	}
	if l.line != "" {
		l.tokens <- token{typ: tNewline}
	}
}

func (l *lexer) parseIndent(r rune) {
	if !unicode.IsSpace(r) {
		cur := l.indent[len(l.indent)-1]
		if l.index > cur {
			l.tokens <- token{typ: tIndent}
			l.indent = append(l.indent, l.index)
		} else if l.index < cur {
			for l.indent[len(l.indent)-1] != l.index {
				l.tokens <- token{typ: tDedent}
				l.indent = l.indent[:len(l.indent)-1]
				if l.index > l.indent[len(l.indent)-1] {
					log.Println("Line #%d: Inconsistent indenting", l.lineNum)
					l.indent[len(l.indent)-1] = l.index
				}
			}
		}
		l.start = l.index
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

func (l *lexer) emit(typ tokenType, v string) {
	l.tokens <- token{typ: typ, value: v}
	l.index += len(v)
	l.start = l.index
}

// The longest sequence creating a valid token is lexed with the following priority:
//	(1) Indent, Dedent, Newline
//	(2)	Keyword, Delimiter
//	(3)	Operator
//	(4)	Identifier, Literal
// Additionally, keywords, identifiers and literals must be separated by whitespace
func (l *lexer) tokenize(r rune, t token) token {
	if t == nil {
		if v := tokenList(keywords, tKeyword, r); v != "" {
			next, _ := l.peek(l.index + len(v))
			if next != '_' && !unicode.IsNumber(next) && !unicode.IsLetter(next) {
				l.emit(tKeyword, v)
				return nil
			}
		}
		if v := tokenList(delimiters, tDelimiter, r); v != "" {
			l.emit(tDelimiter, v)
			return nil
		}
		if v := tokenList(operators, tOperator, r); v != "" {
			l.emit(tDelimiter, v)
			return nil
		}
		if r == '_' || unicode.IsLetter(r) {
			return token{typ: tIdentifier}
		}
	}
}

func (l *lexer) tokenList(list []string, typ tokenType, r rune) string {
	for _, v := range list {
		if strings.HasPrefix(l.line[l.index:], v) {
			return v
		}
	}
	return ""
}

const keywords = []string{
	"byte", "word", "block", "const", "func ", "...", "->", "jump", "return", "if",
}
const delimiters = []string{
	",", "(", ")",
}
const operators = []string{
	":", ":=", "=", "+", "-", "*", "/", "**",
}
