package main

import (
	"io"
	"log"
	"strings"
	"unicode"
)

type parser struct {
	*strings.Reader
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

func (p *parser) getVariable(n uint) *variable {
	return &p.variable[n-1]
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
	p.start.stmt = p.parseStatementList()
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
	// Add special variables
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
		v := p.getWord()
		p.imported = append(p.imported, p.getVariable(v))
	}
}

func (p *parser) parseStatementList() (s []statement) {
	n := p.getWord()
	if n == 0 {
		log.Fatal("Missing statement")
	}
	for ; n > 0; n-- {
		s = append(s, p.parseStatement())
	}
	return
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

func (p *parser) parseStatement() (s statement) {
	switch p.next() {
	case bVariableDef:
		s = p.parseVariableDef()
	case bIf:
		s = p.parseIf()
	case bAssignment:
		s = p.parseAssignment()
	case bJump:
		s = p.parseJump()
	case bReturn:
		s = p.parseReturn()
	default:
		log.Fatal("Invalid statement")
	}
	return
}

func (p *parser) parseVariableDef() statement {
	s := new(variableDefStmt)
	v := p.parseDeclaration()
	p.cur.local = append(p.cur.local, v)
	v.scope = p.cur
	v.base = p.getVariable(p.getWord())
	v.offset = int(p.getWord())
	return s
}

func (p *parser) parseDeclaration() (v *variable) {
	v = p.getVariable(p.getWord())
	switch p.next() {
	case bWord:
		v.refLength = p.wordLength
		v.length = 1
	case bByte:
		v.refLength = 1
		v.length = 1
	case bBlockWord:
		v.refLength = p.wordLength
		v.length = p.getWord()
	case bBlockByte:
		v.refLength = 1
		v.length = p.getWord()
	default:
		log.Fatal("Invalid variable declaration")
	}
	return
}

func (p *parser) parseIf() statement {
	s := new(ifStmt)
	p.parseExpression()
	s.stmt = p.parseStatementList()
	return s
}

func (p *parser) parseAssignment() statement {
	s := new(assignmentStmt)
	for n := p.getWord(); n > 0; n-- {
		v := p.getWord()
		s.assignee = append(s.assignee, p.getVariable(v))
	}
	switch p.next() {
	case bExpression:
		s.expr = p.parseExpression()
	case bFunction:
	}
	return s
}

func (p *parser) parseJump() statement {
	s := new(jumpStmt)
	s.expr = p.parseExpression()
	return s
}

func (p *parser) parseReturn() statement {
	return new(returnStmt)
}

func (p *parser) parseExpression() (e expression) {
	switch p.next() {
	case bLiteral:
		e = p.parseLiteral()
	case bVariableRef:
		e = p.parseVariableRef()
	case bFunctionCall:
		e = p.parseFunctionCall()
	default:
		log.Fatal("Invalid expression")
	}
	return
}

func (p *parser) parseLiteral() expression {
	return
}

func (p *parser) parseVariableRef() expression {
	return
}

func (p *parser) parseFunctionCall() expression {
	return
}

func (p *parser) parseFunction(v *variable) {
	f := &function{bind: v}
	p.cur = f
	for n := p.getWord(); n > 0; n-- {
		p.parseVariableDef()
	}
	p.cur = v.scope
}
