package main

import (
	"bufio"
	"io"
	"strings"
	"strconv"
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
	"byte", "block", "func", "jump", "return", "if", "ref",
}
var operators = [...]string{
	"+", "-", "*", "/", "**", "%", "&", "|", "^", "!", "<<", ">>",
}
var delimiters = [...]string{
	":", ":=", "=", ",", "->", "(", ")",
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
		for state := lexIndent; ; {
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
// Call lexLine after emitting a token
func lexLine(l *lexer) stateFn {
	l.start = l.pos
	if unicode.IsSpace(l.cur) {
		// Ignore whitespace
		return lexLine
	} else if unicode.IsDigit(l.cur) {
		return lexLiteral(l)
	} else if unicode.IsLetter(l.cur) || l.cur == '_' {
		return lexIdentifier(l)
	} else {
		return lexOperator(l)
	}
}

// inList returns a token (without updating *lexer indices) when l.start matches
// a prefix string in the input list
func (l *lexer) inList(t terminal, s []string) *token {
	var tok *token
	for _, k := range(s) {
		if strings.HasPrefix(l.line[l.pos:], k) {
			if tok != nil && len(k) < len(tok.lexeme) {
				continue
			}
			tok = new(token)
			tok.terminal = t
			tok.lexeme = l.line[l.pos:(l.pos + len(k))]
		}
	}
	return tok
}

// lexOperator grabs the longest match in the 'operators' and 'delimiters' lists,
// and emits a token
func lexOperator(l *lexer) stateFn {
	var tok token
	if t := l.inList(tOperator, operators[:]); t != nil {
		tok = *t
	}
	if t := l.inList(tDelimiter, delimiters[:]); t != nil {
		if len(t.lexeme) > len(tok.lexeme) {
			tok = *t
		}
	}
	if tok.lexeme == "" {
		log.Fatal("Invalid operator")
	}
	l.pos += len(tok.lexeme)
	l.width = 0
	l.tokens <- tok
	return lexLine
}

// lexIdentifier emits an identifier token or a reserved token, if the
// identifier is a reserved keyword
func lexIdentifier(l *lexer) stateFn {
	if !unicode.IsLetter(l.cur) && !unicode.IsDigit(l.cur) && l.cur != '_' {
		tok := token{
			terminal:	tIdentifier,
			lexeme:		l.line[l.start:l.pos],
		}
		for _, k := range(reserved) {
			if tok.lexeme == k {
				tok.terminal = tReserved
			}
		}
		if unicode.IsSpace(l.cur) || l.inList(tDelimiter, delimiters[:]) != nil || l.cur == ':' || l.cur == '=' {
			l.tokens <- tok
			return lexLine(l)
		} else {
			log.Fatal("Invalid identifier")
		}
	}
	if l.pos + l.width >= len(l.line) {
		l.next()
		l.cur = ' '
		return lexIdentifier(l)
	}
	return lexIdentifier
}

func lexLiteral(l *lexer) stateFn {
	if !unicode.IsDigit(l.cur) {
		str := l.line[l.start:l.pos]
		tok := token{ terminal: tLiteral }
		var err error
		if l.cur == 'b' {
			tok.num, err = strconv.ParseInt(str, 2, 64)
		} else if l.cur == 'o' {
			tok.num, err = strconv.ParseInt(str, 8, 64)
		} else if l.cur == 'h' {
			tok.num, err = strconv.ParseInt(str, 16, 64)
		} else {
			tok.num, err = strconv.ParseInt(str, 10, 64)
		}
		if err != nil {
			log.Fatal(err)
		}
		l.tokens <- tok
		if (l.cur == 'b' || l.cur == 'o' || l.cur == 'h') && l.next() {
			if !unicode.IsSpace(l.cur) && l.inList(tDelimiter, delimiters[:]) == nil {
				log.Fatal("Invalid literal")
			}
			return lexLine(l)
		}
		return lexLine
	}
	if l.pos + l.width >= len(l.line) {
		l.next()
		l.cur = ' '
		return lexLiteral(l)
	}
	return lexLiteral
}
