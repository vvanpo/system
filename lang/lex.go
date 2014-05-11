package main

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type terminal int

const (
	tNewline terminal = iota // Newlines always end statements and are thus necessary for parsing
	tIndent
	tDedent
	tLiteral
	tIdentifier
	// Fixed lexemes (found in lexer struct)
	tByte
	tBlock
	tFunc
	tJump
	tReturn
	tIf
	tRef
	tAdd
	tSub
	tMult
	tDiv
	tExp
	tMod
	tAnd
	tOr
	tXor
	tNot
	tShiftL
	tShiftR
	tAlias
	tAutoVar
	tAssign
	tComma
	tMap
	tLeftParen
	tRightParen
)

type token struct {
	terminal
	lexeme string
	line   int // line & col information passed to parser for anotating errors
	col    int
}

type lexer struct {
	line        string // Current line
	lineNum     int
	tokens      chan token
	indent      []int
	indent_rune rune                // Tracks the rune used for indentation
	start       int                 // Start of current token
	pos         int                 // Start of current rune
	cur         rune                // Rune pointed to by pos
	width       int                 // Width of cur
	reserved    map[string]terminal // Map of reserved words to terminals
	key         map[string]terminal // Map of all other tokens indexed by fixed lexemes
}

type lexFn func(*lexer) lexFn

func (l *lexer) lexErr(s string) {
	log.Printf("Ln %d, col %d: lexing error", l.lineNum, l.pos+1)
	if s != "" {
		log.Fatalf(":\n\t%s", s)
	} else {
		log.Fatal("\n")
	}
}

// lex initializes goroutine to lex input and returns token channel
// When channel closes, we've reached EOF
func lex(r io.Reader) chan token {
	l := &lexer{
		tokens: make(chan token),
		indent: []int{0},
		reserved: map[string]terminal{
			"byte":   tByte,
			"block":  tBlock,
			"func":   tFunc,
			"jump":   tJump,
			"return": tReturn,
			"if":     tIf,
			"ref":    tRef,
		},
		key: map[string]terminal{
			":=": tAutoVar,
			"=":  tAssign,
			":":  tAlias,
			",":  tComma,
			"->": tMap,
			"(":  tLeftParen,
			")":  tRightParen,
			"+":  tAdd,
			"-":  tSub,
			"*":  tMult,
			"/":  tDiv,
			"**": tExp,
			"%":  tMod,
			"&":  tAnd,
			"|":  tOr,
			"^":  tXor,
			"!":  tNot,
			"<<": tShiftL,
			">>": tShiftR,
		},
	}
	go l.run(bufio.NewScanner(r))
	return l.tokens
}

// *lexer.run lexes entire file line-by-line
func (l *lexer) run(s *bufio.Scanner) {
	defer close(l.tokens)
	for s.Scan() {
		l.line = s.Text()
		l.lineNum++
		l.start = 0
		l.pos = 0
		l.width = 0
		for state := lexIndent; ; {
			if !l.next() {
				l.emit(token{
					tNewline, "", l.lineNum, len(l.line),
				})
				break
			}
			state = state(l)
		}
	}
	if s.Err() != nil {
		l.lexErr("Input error")
	}
}

func (l *lexer) emit(t token) {
	t.line = l.lineNum
	t.col = l.start + 1
	l.tokens <- t
}

// isLast returns true if l.cur is the last rune in the line
func (l *lexer) isLast() bool {
	if l.pos+l.width == len(l.line) || l.line[l.pos+l.width] == uint8('#') {
		return true
	}
	return false
}

// *lexer.next updates *lexer values and returns true if there is a valid rune to lex
func (l *lexer) next() bool {
	if l.isLast() {
		return false
	}
	l.pos += l.width
	r, s := utf8.DecodeRuneInString(l.line[l.pos:])
	l.width = s
	if r == '\ufffd' && s == 1 {
		l.lexErr("Invalid input encoding")
	}
	l.cur = r
	return true
}

// Call lexIndent at the beginning of a line
func lexIndent(l *lexer) lexFn {
	if !unicode.IsSpace(l.cur) {
		for i := range l.indent {
			if l.indent[i] == l.pos {
				if diff := len(l.indent) - 1 - i; diff != 0 {
					l.indent = l.indent[:i+1]
					for j := 0; j < diff; j++ {
						l.emit(token{terminal: tDedent})
					}
				}
				return lexNext(l)
			} else if l.indent[i] > l.pos {
				l.lexErr("Indentation mismatch")
			}
		}
		l.indent = append(l.indent, l.pos)
		l.emit(token{terminal: tIndent})
		return lexNext(l)
	}
	if l.cur != '\t' && l.cur != ' ' {
		l.lexErr("Invalid indentation character")
	} else if l.indent_rune == 0 {
		l.indent_rune = l.cur
	} else if l.cur != l.indent_rune {
		l.lexErr("Mixing tabs and spaces for indentation")
	}
	return lexIndent
}

// lexNext determines which state function to call for current rune
// Call lexNext after emitting a token
func lexNext(l *lexer) lexFn {
	l.start = l.pos
	if unicode.IsSpace(l.cur) {
		// Ignore whitespace
		return lexNext
	} else if unicode.IsDigit(l.cur) {
		return lexLiteral(l)
	} else if unicode.IsLetter(l.cur) || l.cur == '_' {
		return lexIdentifier(l)
	} else {
		return lexFixed(l)
	}
}

// lexLiteral emits tokens after converting values of the forms:
//		<binary>b
//		<octal>o
//		<hexadecimal>h
// to decimal form
func lexLiteral(l *lexer) lexFn {
	parseInt := func(i, sz int) {
		var base int
		switch string(l.line[i]) {
		case "b":
			base = 2
		case "o":
			base = 8
		case "h":
			base = 16
		default:
			base = 10
			i += sz
		}
		n, err := strconv.ParseInt(l.line[l.start:i], base, 64)
		if err != nil {
			l.lexErr(err.Error())
		}
		t := token{terminal: tLiteral}
		t.lexeme = strconv.Itoa(int(n))
		l.emit(t)
	}
	if l.cur == ':' || l.cur == '=' || unicode.IsSpace(l.cur) {
		_, sz := utf8.DecodeLastRuneInString(l.line[:l.pos])
		parseInt(l.pos-sz, sz)
		return lexNext(l)
	} else if l.isLast() {
		parseInt(l.pos, l.width)
		return lexNext
	}
	return lexLiteral
}

// lexIdentifier emits an identifier token or a reserved token, if the
// identifier is a reserved keyword
func lexIdentifier(l *lexer) lexFn {
	parseIdentifier := func(i int) {
		t := token{terminal: tIdentifier, lexeme: l.line[l.start:i]}
		s := l.line[l.start:i]
		for k := range l.reserved {
			if s == k {
				t.terminal = l.reserved[k]
			}
		}
		l.emit(t)
	}
	if unicode.In(l.cur, unicode.Nd, unicode.L) || l.cur == '_' {
		if l.isLast() {
			parseIdentifier(l.pos + l.width)
		}
		return lexIdentifier
	} else if l.cur == ',' || l.cur == ')' || l.cur == ':' || l.cur == '=' || unicode.IsSpace(l.cur) {
		parseIdentifier(l.pos)
		return lexNext(l)
	}
	l.lexErr("Invalid identifier")
	return nil
}

func lexFixed(l *lexer) lexFn {
	parseFixed := func(i int) terminal {
		s := l.line[l.start:i]
		var t token
		for k := range l.key {
			if s == k {
				t = token{terminal: l.key[k]}
			}
		}
		if t.terminal == 0 {
			l.lexErr("Invalid symbol")
		}
		l.emit(t)
		return t.terminal
	}
	if unicode.In(l.cur, unicode.Nd, unicode.L, unicode.Z) || l.cur == '_' {
		t := parseFixed(l.pos)
		if !unicode.IsSpace(l.cur) {
			switch t {
			default: l.lexErr("No whitespace after symbol '" + l.line[l.start:l.pos] + "'")
			case tAlias:
			case tAutoVar:
			case tAssign:
			case tComma:
			case tLeftParen:
			}
		}
		return lexNext(l)
	} else if l.isLast() {
		parseFixed(l.pos + l.width)
		return nil
	}
	return lexFixed
}
