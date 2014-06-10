package main

import (
	"bufio"
	"io"
	"log"
	"unicode"
)

type parser struct {
	*bufio.Reader
	bytelang
}

func newParser(r io.Reader) (p *parser) {
	p = &parser{Reader: bufio.NewReader(r)}
	return
}

func (p *parser) next() (c byte) {
	c, err := p.ReadByte()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (p *parser) read(n int) (s string) {
	b := make([]byte, n)
	if _, err := p.Read(b); err != nil {
		log.Fatal(err)
	}
	s = string(b)
	return
}

func (p *parser) parseWord() (word uint) {
	b := p.read(p.wordLength)
	for i := 0; i < len(b); i++ {
		word |= uint(b[i])
		word <<= 8
	}
	return
}

func (p *parser) addIdentifier(id string) {
	for i, r := range id {
		if r == '_' || unicode.IsLetter(r) || (unicode.IsDigit(r) && i != 0) {
			continue
		}
		log.Fatal("Invalid identifier")
	}
	for _, s := range p.identifier {
		if s == identifier(id) {
			log.Println("Duplicate identifier")
		}
	}
	p.identifier = append(p.identifier, identifier(id))
}

func (p *parser) addStatement(s statement) {
	p.statement = append(p.statement, s)
}

func (p *parser) parseBytelang() {
	p.parseHeader()
	p.parseIdentifierList()
	p.parseStatementList()
}

func (p *parser) parseHeader() {
	match := "Version 0.0\nArch.: "
	if match != p.read(len(match)) {
		log.Fatal("Invalid header string")
	}
	c, err := p.ReadString(' ')
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(c)-2; i++ {
		if c[i] < byte('0') || c[i] > byte('9') {
			log.Fatal("Invalid word-length definition")
		}
		p.wordLength = (p.wordLength * 10) + int(c[i]-byte('0'))
	}
	match = "bytes/word\n"
	if match != p.read(len(match)) {
		log.Fatal("Invalid header string")
	}
}

func (p *parser) parseIdentifierList() {
	n := p.parseWord()
	if specialIdentifiers != p.read(len(specialIdentifiers)) {
		log.Fatal("Missing special identifier")
	}
	last := 0
	for i, r := range specialIdentifiers {
		if r == '\n' {
			p.addIdentifier(specialIdentifiers[last:i])
			last = i
			n--
		}
	}
	for ; n > 0; n-- {
		id, err := p.ReadString('\n')
		if err != nil {
			log.Fatal("Invalid identifier list")
		}
		p.addIdentifier(id)
	}
}

func (p *parser) parseStatementList() {
	n := p.parseWord()
	if n == 0 {
		log.Fatal("Missing statement")
	}
	for ; n > 0; n-- {
		p.parseStatement()
	}
	return
}

func (p *parser) parseStatement() {
	switch p.next() {
	case bVariable:
		p.parseVariable()
	case bFunctionCall:
		p.parseFunctionCall()
	case bIf:
		p.parseIf()
	case bAssignment:
		p.parseAssignment()
	case bReturn:
		p.parseReturn()
	default:
		log.Fatal("Invalid statement")
	}
	return
}

func (p *parser) parseVariable() {
	i := p.parseWord()
	switch p.next() {
	case bWord:
		v := &variableWord{
			bytelang:   &p.bytelang,
			identifier: &p.identifier[i-1],
		}
		p.addStatement(v)
	case bByte:
		v := &variableByte{
			bytelang:   &p.bytelang,
			identifier: &p.identifier[i-1],
		}
		p.addStatement(v)
	case bBlock:
		v := &variableBlock{length: p.parseWord()}
		v.bytelang = &p.bytelang
		v.identifier = &p.identifier[i-1]
		for i := p.parseWord(); i > 0; i-- {
			m := struct {
				offset uint
				v      statement
			}{
				p.parseWord(),
				p.statement[p.parseWord()],
			}
			v.member = append(v.member, m)
		}
		p.addStatement(v)
	default:
		log.Fatal("Invalid variable definition")
	}
}

func (p *parser) parseFunctionCall() {

}

func (p *parser) parseIf() {
	e := p.parseExpression()
}

func (p *parser) parseAssignment() {
	e := p.parseExpression()
}

func (p *parser) parseReturn() {
}

func (p *parser) parseExpression() (e expression) {
	switch p.next() {
	default:
		log.Fatal("Invalid expression")
	}
	return
}
