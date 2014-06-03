package main

import (
	"bytes"
	"io"
	"log"
	"unicode"
)

type tree struct {
	child []*tree
}

func (t *tree) add(child *tree) {
	t.child = append(t.child, child)
}

type parser struct {
	*bytes.Buffer
	fileLen    int64
	wordLen    int // Bytes per word
	syntaxTree *tree
	identifier []identifier
	symbol     []symbol
	imported   []*symbol
	literal    []literal
}

func newParser(r io.Reader) (p *parser) {
	p = new(parser)
	p.Buffer = new(bytes.Buffer)
	n, err := p.ReadFrom(r)
	if err != nil {
		log.Fatal(err)
	}
	p.fileLen = n
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

func (p *parser) addSymbol(idNumber int, address uint) {
	s := symbol{
		&p.identifier[idNumber-1],
		address,
		nil,
	}
	p.symbol = append(p.symbol, s)
}

func (p *parser) getWord() (word uint) {
	b := p.Next(p.wordLen)
	if len(b) < p.wordLen {
		log.Fatal("Not enough bytes in datastream to fill word")
	}
	for i := range b {
		word |= uint(b[i]) << (8 * uint(p.wordLen-i-1))
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

func (p *parser) parseFile() {
	p.parseHeader()
	p.parseIdentifierList()
	p.parseSymbolTable()
	p.parseImportTable()
	p.parseStartBytecode()
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
		p.wordLen = (p.wordLen * 10) + int(c[i]-byte('0'))
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

func (p *parser) parseSymbolTable() {
	n := p.getWord()
	for i := range p.identifier {
		p.addSymbol(i+1, 0)
		n--
	}
	for ; n > 0; n-- {
		id := p.getWord()
		address := p.getWord()
		p.addSymbol(int(id), address)
	}
}

func (p *parser) parseImportTable() {
	for n := p.getWord(); n > 0; n-- {
		symbol := p.getWord()
		p.imported = append(p.imported, &p.symbol[symbol-1])
	}
}

func (p *parser) parseStartBytecode() {
	n := p.getWord()
	if n == 0 {
		log.Fatal("Missing definition statement")
	}
	p.syntaxTree = new(tree)
	for ; n > 0; n-- {
		p.parseSymbolDef()
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

func (p *parser) parseSymbolDef() {
	if p.next() != bSymbolDef {
		log.Fatal("Invalid definition statement")
	}
	n := p.getWord()
	switch p.next() {
	case bAutomatic:
		p.parseDeclaration()
	case bAddress:
		p.parseDeclaration()
		p.parseExpression()
	case bOffset:
		p.parseOffset()
	default:
		log.Fatal("Invalid definition statement")
	}
}
