
package main

import (

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

func (p *parser) next() bool {
	t, ok := <-p.tokens
	if !ok { return false }
	p.list = append(p.list, t)
	return true
}

func (n *node) addChild(c *node) {
	c.parent = n
	n.child = append(n.child, c)
	return
}

func (p *parser) run() {
	for p.parseStmt() {
		p.tree.addChild(p.cur)
	}
}

func (p *parser) parseStmt() bool {
	root := new(node)
	defer func() { p.cur = root }()
	if p.parseLabel() {
		root.addChild(p.cur)
		// The parse<stmt> methods assume p.cur is the stmt root
		p.cur = root
		if p.parseFunc_def() || p.parseIf_stmt() || p.parseBlock() {
			root.addChild(p.cur)
			return true
		}
		// In the case of label_stmt, the owning non-terminal is ambiguous
		if p.parseExpr() {
			root.addChild(p.cur)
			root.nonterm = nLabel_stmt
			return true
		}
	}
	if p.parseAssign_stmt() || p.parseLabel_stmt() || p.parseReassign_stmt() || p.parseJump_stmt() || p.parseReturn_stmt() {
		root.addChild(p.cur)
		return true
	}
	log.Fatal("Invalid statement")
}





/*
func (p *parser) parseStmt() {
	p.addChild()
	tList := make([]token, 5)
	tList[0] = p.curToken
	update := func() {
		p.nextToken()
		tList = append(tList, p.curToken)
	}
	pop := func(int i) {
			tList = tList[(i - 1):]
		}
	// Parse for label
	if tList[0].terminal == tIdentifier {
		isLabel := true
		update()
		// Parse param
		if tList[0].lexeme == "byte" || tList[0].lexeme == "block" {
			p.addChild()
			p.cur.nonterm = nType
			// 0 indicates "byte"
			p.cur.value = 0
			if tList[0].lexeme == "block" {
				if tList[1].terminal != tLiteral {
					log.Fatal("Invalid block size")
				}
				p.cur.value = tList[1].num
				pop(1)
			}
			pop(1)
			p.cur = p.cur.parent
			isLabel = false
		}
		if pList[1].lexeme == ":" {
			if !isLabel {
				
			}
		}
	}

	if curToken.terminal == tReserved {
		if curToken.lexeme == "jump" {
			p.cur.nonterm = nJump_stmt
			p.nextToken()
			p.parseExpr()
			if len(p.cur.child) != 1 {
				log.Fatal("Invalid jump statement")
			}
		} else if curToken.lexeme == "return" {
			p.cur.nonterm = nReturn_stmt
		}
	} else if curToken.terminal == tIdentifier {
		token := curToken
		p.nextToken()
		if p.curToken.lexeme == ":" {
			p.nextToken()
			p.parseExpr()
			if len(p.cur.child) > 0 {
				p.cur.nonterm = nLabel_stmt
			}
		} else 
	} else {
		log.Fatal("Invalid statement")
	}
	p.cur = p.cur.parent
}

func (p *parser) parseExpr() {
	p.addChild()
	p.cur = p.cur.parent
}

*/
