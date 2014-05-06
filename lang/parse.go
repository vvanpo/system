
package main

import (

)

type nonterm int

const (
	stmt nonterm = iota
	func_def
)

type parser struct{
	tokens	chan token
	tree	*node
	cur		*node
}

type node struct{
	parent	*node
	child	[]node
	token
}

func parse(tokens chan token) {
	p := parser{
		tokens:	tokens,
		tree:	new(node),
	}
	p.tree.parent = p.tree
	go p.run()
	return
}

func (p *parser) run() {
	for p.parseStmt() {}
}

func (p *parser) parseStmt() bool {
	return false
}
