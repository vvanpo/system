
package main

import (
	"log"
)

type nonterm int

const (
	nLabel nonterm = iota

	nFunc_def
	nBlock
	nIf_stmt
	nAssign_stmt
	nLabel_stmt
	nReassign_stmt
	nJump_stmt
	nReturn_stmt

	nParam
	nType
	nExpr
	nFunc_call
)

type parser struct{
	tokens	chan token
	list	[]token
	tree	*node
	cur		*node
}

type node struct{
	parent	*node
	child	[]*node
	nonterm
	value	int64
	symbol	string
}

func parse(tokens chan token) {
	p := parser{
		tokens:	tokens,
		tree:	new(node),
	}
	p.tree.parent = p.tree
	p.cur = p.tree
	go p.parseFile()
	return
}

func (p *parser) nextToken() bool {
	t, ok := <-p.tokens
	if !ok { return false }
	p.list = append(p.list, t)
	return true
}

func (p *parser) getToken(i int) (token, bool) {
	if d := i - len(p.list); d >= 0 {
		for j := 0; j <= d; j++ {
			if !p.nextToken() {
				return token{}, false
			}
		}
	}
	return p.list[i], true
}

func (n *node) addChild(c *node) {
	c.parent = n
	n.child = append(n.child, c)
	return
}

func (p *parser) parseFile() bool {
	root := p.cur
	if !p.parseStmt() { return false }
	root.addChild(p.cur)
	for p.parseStmt() {
		root.addChild(p.cur)
	}
	p.cur = root
	return true
}

func (p *parser) parseStmt() bool {
	saved := p.cur
	root := new(node)
	if p.parseLabel() {
		root.addChild(p.cur)
		// The parse<stmt> methods assume p.cur is the stmt root
		p.cur = root
		if p.parseFunc_def() || p.parseIf_stmt() || p.parseBlock() {
			p.cur = root
			return true
		}
		// In the case of label_stmt, the owning non-terminal is ambiguous
		if p.parseExpr() {
			label := root.child[0]
			root.child = root.child[0:0]
			root.nonterm = nLabel_stmt
			// Turn label identifier into param
			param := &node{ nonterm: nParam }
			param.addChild(label)
			root.addChild(param)
			// Add expression to label_stmt
			root.addChild(p.cur)
			p.cur = root
			return true
		}
	}
	if p.parseAssign_stmt() || p.parseLabel_stmt() || p.parseReassign_stmt() || p.parseJump_stmt() || p.parseReturn_stmt() {
		p.cur = root
		return true
	}
	// If the statement was invalid, check if a label was parsed
	if len(root.child) > 0 {
		log.Fatal("Invalid statement")
	}
	p.cur = saved
	return false
}

func (p *parser) parseLabel() bool {
	if t, ok := p.getToken(0); ok && t.terminal == tIdentifier {
		if u, ok := p.getToken(1); ok && u.lexeme == ":" {
			p.cur = &node{
				parent:		p.cur,
				nonterm:	nLabel,
				symbol:		t.lexeme,
			}
			p.list = p.list[2:]
			return true
		}
	}
	return false
}

func (p *parser) parseFunc_def() bool {
	root := p.cur
	if t, ok := p.getToken(0); !ok || t.lexeme != "func" { return false }
	p.list = p.list[1:]
	for p.parseParam() {
		root.addChild(p.cur)
	}
	if t, ok := p.getToken(0); ok && t.lexeme == "->" {
		p.list = p.list[1:]
		if !p.parseParam() {
			log.Fatal("Invalid function definition")
		}
		root.addChild(p.cur)
		for p.parseParam() {
			root.addChild(p.cur)
		}
	}
	// Block node
	p.cur = new(node)
	root.addChild(p.cur)
	if !p.parseBlock() {
		log.Fatal("Invalid function definition")
	}
	root.nonterm = nFunc_def
	p.cur = root
	return true
}

func (p *parser) parseBlock() bool {
	if t, ok := p.getToken(0); !ok || t.terminal != tIndent {
		return false
	}
	p.list = p.list[1:]
	p.cur.nonterm = nBlock
	if !p.parseFile() {
		log.Fatal("Invalid block statement")
	}
	if t, ok := p.getToken(0); !ok || t.terminal != tDedent {
		log.Fatal("Invalid block statement")
	}
	p.list = p.list[1:]
	return true
}

func (p *parser) parseIf_stmt() bool {
	if t, ok := p.getToken(0); !ok || t.lexeme != "if" {
		return false
	}
	p.list = p.list[1:]
	root := p.cur
	root.nonterm = nIf_stmt
	if !p.parseExpr() {
		log.Fatal("Invalid if-statement")
	}
	root.addChild(p.cur)
	// Block node
	p.cur = new(node)
	root.addChild(p.cur)
	if !p.parseBlock() {
		log.Fatal("Invalid function definition")
	}
	p.cur = root
	return true
}

func (p *parser) parseParam() bool {
	root := new(node)
	root.nonterm = nParam
	parseSingle := func() bool {
		t, ok := p.getToken(0)
		if !ok || ( t.terminal != tReserved && t.terminal != tIdentifier ) {
			return false
		}
		if t.terminal == tReserved {
			if t.lexeme == "byte" {
				root.addChild(&node{ nonterm: nType, value: 0 })
				p.list = p.list[1:]
			} else if t.lexeme == "block" {
				if t, ok := p.getToken(1); ok && t.terminal == tLiteral {
					root.addChild(&node{ nonterm: nType, value: t.num })
					p.list = p.list[2:]
				} else {
					log.Fatal("Invalid block-typed parameter")
				}
			}
			t, ok = p.getToken(0)
			if !ok || t.terminal != tIdentifier {
				log.Fatal("Invalid parameter")
			}
		}
		root.addChild(&node{ symbol: t.lexeme })
		p.list = p.list[1:]
		return true
	}
	if !parseSingle() { return false }
	for parseSingle() {}
	p.cur = root
	return true
}

