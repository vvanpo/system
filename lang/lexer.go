package main

import (
	"bufio"
	"io"
	//"strings"
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
	"byte", "block", "func", "jump", "return", "if",
}
var operators = [...]string{
	":", ":=", "=", "+", "-", "*", "/", "**", "(", ")",
}
var delimiters = [...]string{
	",", "->",
}

type token struct {
	terminal
	lexeme string
	num    int64
}

type stateFn func(*lexer) stateFn

type lexer struct {
	line        string // Current line
	tokens      chan token
	indent      []int
	indent_rune rune
	start       int  // Start of current token
	pos         int  // Start of current rune
	cur         rune // Rune pointed to by pos
	width       int  // Width of cur
}

// lex initializes goroutine to lex input and returns token channel
// When channel closes, we've reached EOF
func lex(r io.Reader) chan token {
	l := &lexer{
		tokens: make(chan token),
		indent: []int{0},
	}
	go l.run(bufio.NewScanner(r))
	return l.tokens
}

// *lexer.run lexes entire file line-by-line
func (l *lexer) run(s *bufio.Scanner) {
	defer close(l.tokens)
	for s.Scan() {
		l.line = s.Text()
		l.start = 0
		l.pos = 0
		l.width = 0
		for state := lexIndent; state != nil; {
			if !l.next() {
				break
			}
			state = state(l)
		}
	}
	if s.Err() != nil {
		log.Fatal("Input error")
	}
}

// *lexer.next updates *lexer values and returns true if there is a valid rune to lex
func (l *lexer) next() bool {
	l.pos += l.width
	if l.pos >= len(l.line) {
		return false
	}
	r, s := utf8.DecodeRuneInString(l.line[l.pos:])
	l.width = s
	if r == '#' {
		return false
	}
	if r == '\ufffd' && s == 1 {
		log.Fatal("Invalid input encoding")
	}
	l.cur = r
	return true
}

// Call lexIndent at the beginning of a line
func lexIndent(l *lexer) stateFn {
	if !unicode.IsSpace(l.cur) {
		l.start = l.pos
		for i := range l.indent {
			if l.indent[i] == l.pos {
				if diff := len(l.indent) - 1 - i; diff != 0 {
					l.indent = l.indent[:i+1]
					for j := 0; j < diff; j++ {
						l.tokens <- token{terminal: tDedent}
					}
				}
				return lexLine(l)
			} else if l.indent[i] > l.pos {
				log.Fatal("Indentation mismatch")
			}
		}
		l.indent = append(l.indent, l.pos)
		l.tokens <- token{terminal: tIndent}
		return lexLine(l)
	}
	if l.cur != '\t' && l.cur != ' ' {
		log.Fatal("Invalid indentation character")
	} else if l.indent_rune == 0 {
		l.indent_rune = l.cur
	} else if l.cur != l.indent_rune {
		log.Fatal("Mixing tabs and spaces for indentation")
	}
	return lexIndent
}

// lexLine determines which state function to call for current rune
func lexLine(l *lexer) stateFn {
	return nil
}

/*

// lexIdentifier emits an identifier token or a reserved token, if the
// identifier is a reserved keyword
func lexIdentifier(l *lexer) stateFn {
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


*/
