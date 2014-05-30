package main

import (
	"io"
	"regexp"
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
}

func newParser(r io.Reader) (p *parser) {
	p = new(parser)
	p.Reader = r
}

func (p *parser) parseFile() {
	p.parseHeader()
	//	p.parseIdentifierList()
	//	p.parseSymbolTable()
	//	p.parseImportTable()
	//	p.parseStartBytecode()
	//	p.parseLiteralList()
}

func (p *parser) parseHeader() {

}
