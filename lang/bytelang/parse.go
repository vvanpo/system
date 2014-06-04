package main

import (
	"bytes"
	"io"
	"log"
	"unicode"
)

type parser struct {
	*bytes.Buffer
	bytelang
	cur *function // Function scope tracking during parsing
}

func newParser(r io.Reader) (p *parser) {
	_, err := p.ReadFrom(r)
	if err != nil {
		log.Fatal(err)
	}
	p = &parser{
		Buffer: new(bytes.Buffer),
	}
	p.cur = &p.start
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

func (p *parser) addVariable(idNumber int) {
	s := variable{
		identifier: &p.identifier[idNumber-1],
	}
	p.variable = append(p.variable, s)
}

func (p *parser) getWord() (word uint) {
	b := p.Next(p.wordLength)
	if len(b) < p.wordLength {
		log.Fatal("Not enough bytes in datastream to fill word")
	}
	for i := range b {
		word |= uint(b[i]) << (8 * uint(p.wordLength-i-1))
	}
	return
}

func (p *parser) next() (c byte) {
	c, err := p.ReadByte()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (p *parser) parseBytelang() {
	p.parseHeader()
	p.parseIdentifierList()
	p.parseVariableTable()
	p.parseImportTable()
	p.parseStatementList()
	p.parseLiteralList()
}

func (p *parser) parseHeader() {
	match := "Version 0.0\nArch.: "
	if match != string(p.Next(len(match))) {
		log.Fatal("Invalid header string")
	}
	c, err := p.ReadBytes(byte(' '))
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	for i := 0; i < len(c)-2; i++ {
		if c[i] < byte('0') || c[i] > byte('9') {
			break
		}
		p.wordLength = (p.wordLength * 10) + int(c[i]-byte('0'))
	}
	match = "bytes/word\n"
	if match != string(p.Next(len(match))) {
		log.Fatal("Invalid header string")
	}
}

func (p *parser) parseIdentifierList() {
	n := p.getWord()
	if specialIdentifiers != string(p.Next(len(specialIdentifiers))) {
		log.Fatal("Missing special identifiers")
	}
	last := 0
	for i, r := range specialIdentifiers {
		if r == 0xfffd {
			log.Fatal("Invalid identifier encoding")
		}
		if r == '\n' {
			p.addIdentifier(specialIdentifiers[last:i])
			last = i
			n--
		}
	}
	for ; n > 0; n-- {
		id, err := p.ReadString(byte('\n'))
		if err != nil {
			log.Fatal("Invalid identifier list")
		}
		p.addIdentifier(id)
	}
}

func (p *parser) parseVariableTable() {
	n := p.getWord()
	for i := range p.identifier {
		p.addVariable(i + 1)
		n--
	}
	for ; n > 0; n-- {
		id := p.getWord()
		p.addVariable(int(id))
	}
}

func (p *parser) parseImportTable() {
	for n := p.getWord(); n > 0; n-- {
		variable := p.getWord()
		p.imported = append(p.imported, &p.variable[variable-1])
	}
}

func (p *parser) parseStatementList() {
	n := p.getWord()
	if n == 0 {
		log.Fatal("Missing statement")
	}
	for ; n > 0; n-- {
		p.parseStatement()
	}
}

func (p *parser) parseLiteralList() {
	for p.Len() > 0 {
		n := p.getWord()
		lit := make(literal, n)
		p.literal = append(p.literal, lit)
		for i := 0; i < int(n); i++ {
			lit[i] = p.next()
		}
	}
}

func (p *parser) parseStatement() {
	switch p.next() {
	case bAutomatic:
		p.parseAutomatic()
	case bAddress:
		p.parseAddress()
	case bOffset:
		p.parseOffset()
	case bIf:
		p.parseIf()
	case bAssignment:
		p.parseAssignment()
	case bJump:
		p.parseJump()
	case bReturn:
		p.parseReturn()
	default:
		log.Fatal("Invalid statement")
	}
}

func (p *parser) parseAutomatic() {
	v := p.parseDeclaration()
}

func (p *parser) parseAddress() {
	v := p.parseDeclaration()
	p.parseExpression()
}

func (p *parser) parseOffset() {
	v := p.parseDeclaration()
	n := p.getWord()      // Parent variable number
	offset := p.getWord() //  Offset index into parent variable
}

func (p *parser) parseDeclaration() (v *variable) {
	n := p.getWord() // Variable number
	v = &p.variable[n-1]
	switch p.next() {
	case bWord:
		v.refLength = p.wordLength
		v.length = uint(p.wordLength)
	case bByte:
		v.refLength = 1
		v.length = 1
	case bBlockWord:
		v.refLength = p.wordLength
		v.length = p.getWord()
	case bBlockByte:
		v.refLength = 1
		v.length = p.getWord()
	case bFunction:
		v.refLength = p.wordLength
		v.length = uint(p.wordLength)
		p.parseFunction(v)
	}
	return
}

func (p *parser) parseIf() {
	p.parseExpression()
	p.parseStatementList()
}

func (p *parser) parseAssignment() {
	for n := p.getWord(); n > 0; n-- {
		v := p.getWord()
	}
	p.parseExpression()
}

func (p *parser) parseJump() {
	p.parseExpression()
}

func (p *parser) parseReturn() {
}

func (p *parser) parseExpression() {
	switch p.next() {
	case bLiteral:
		p.parseLiteral()
	case bVariableRef:
		p.parseVariableRef()
	case bFunctionCall:
		p.parseFunctionCall()
	default:
		log.Fatal("Invalid expression")
	}
}

func (p *parser) parseLiteral() {

}

func (p *parser) parseVariableRef() {

}

func (p *parser) parseFunctionCall() {

}
