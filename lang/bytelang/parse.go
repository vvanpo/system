package main

import (
	"io"
	"log"
)

type tree struct {
	child []*tree
}

func (t *tree) add(child *tree) {
	t.child = append(t.child, child)
}

type parser struct {
	io.Reader
	wordSize   int // Bytes per word
	syntaxTree *tree
	identifier []identifier
	symbol     []symbol
	imported   []symbol
	literal    []literal
}

func newParser(r io.Reader) (p *parser) {
	p = new(parser)
	p.Reader = r
	return
}

func (p *parser) getWord() (word uint) {
	b := make([]byte, p.wordSize)
	n, _ := p.Reader.Read(b)
	if n != p.wordSize {
		log.Fatal("Parsing error")
	}
	for i := range(b) {
		word <<= 8 * uint(i)
		word += uint(b[i])
	}
	return
}

func (p *parser) parseFile() {
	p.parseHeader()
	p.parseIdentifierList()
	//	p.parseSymbolTable()
	//	p.parseImportTable()
	//	p.parseStartBytecode()
	//	p.parseLiteralList()
}

func (p *parser) parseHeader() {
	match := "Version 0.0\nArch.: "
	buf := make([]byte, len(match))
	p.Reader.Read(buf)
	if match != string(buf) {
		log.Fatal("Invalid header string")
	}
	buf = make([]byte, 3)
	for i := 0; i < len(buf); i++ {
		n, _ := p.Reader.Read(buf[i : i+1])
		if n != 1 || buf[i] < byte('0') || buf[i] > byte('9') {
			break
		}
		p.wordSize = (p.wordSize * 10) + int(buf[i]-byte('0'))
	}
	match = "bytes/word\n"
	buf = make([]byte, len(match))
	p.Reader.Read(buf)
	if match != string(buf) {
		log.Fatal("Invalid header string")
	}
}

func (p *parser) parseIdentifierList() {
	n := p.getWord()
}
