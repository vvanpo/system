
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
	for {
		if t, ok := p.getToken(0); ok && t.lexeme == "," {
			p.list = p.list[1:]
			if !parseSingle() {
				log.Fatal("Invalid parameter list")
			}
		} else {
			break
		}
	}
	p.cur = root
	return true
}

func (p *parser) parseAssign_stmt() bool {
	root := p.cur
	defer func() { p.cur = root }()
	if !p.parseParam() { return false }
	param := p.cur
	if t, ok := p.getToken(0); !ok || t.lexeme != ":=" { return false }
	if !p.parseExpr() { return false }
	p.list = p.list[1:]
	root.addChild(param)
	root.addChild(p.cur)
	root.nonterm = nAssign_stmt
	return true
}

func (p *parser) parseLabel_stmt() bool {
	root := p.cur
	defer func() { p.cur = root }()
	if !p.parseParam() { return false }
	param := p.cur
	if t, ok := p.getToken(0); !ok || t.lexeme != ":" { return false }
	if !p.parseExpr() { return false }
	p.list = p.list[1:]
	root.addChild(param)
	root.addChild(p.cur)
	root.nonterm = nLabel_stmt
	return true
}


func (p *parser) parseReassign_stmt() bool {
	root := p.cur
	defer func() { p.cur = root }()
	parseSingle := func() bool {
		t, ok := p.getToken(0)
		if !ok || t.terminal != tIdentifier { return false }
		root.addChild(&node{ symbol: t.lexeme })
		p.list = p.list[1:]
		return true
	}
	if !parseSingle() { return false }
	for {
		if t, ok := p.getToken(0); ok && t.lexeme == "," {
			p.list = p.list[1:]
			if !parseSingle() {
				log.Fatal("Invalid parameter list")
			}
		} else {
			break
		}
	}
	if t, ok := p.getToken(0); !ok || t.lexeme != "=" {
		log.Fatal("Invalid assignment statement")
	}
	p.list = p.list[1:]
	if !p.parseExpr() {
		log.Fatal("Invalid assignment statement")
	}
	root.addChild(p.cur)
	root.nonterm = nReassign_stmt
	return true
}

func (p *parser) parseJump_stmt() bool {
	if t, ok := p.getToken(0); !ok || t.lexeme != "jump" { return false }
	p.list = p.list[1:]
	root := p.cur
	if !p.parseExpr() {
		log.Fatal("Invalid jump statement")
	}
	root.addChild(p.cur)
	root.nonterm = nJump_stmt
	p.cur = root
	return true
}

func (p *parser) parseReturn_stmt() bool {
	if t, ok := p.getToken(0); !ok || t.lexeme != "return" { return false }
	p.list = p.list[1:]
	p.cur.nonterm = nReturn_stmt
	return true
}

func (p *parser) parseExpr() bool {
	nodes := make([]*node, 1)
	subExpr := -1
	i := 0
	for {
		t, ok := p.getToken(i)
		if !ok { break }
		if t.lexeme == "(" {
			tmp := p.list[:i]
			p.list = p.list[i:]
			p.parseFunc_call()
			nodes = append(nodes, p.cur)
			p.list = append(tmp, p.list...)
			continue
		}
		if t.terminal == tIdentifier {
			nodes = append(nodes, &node{ symbol: t.lexeme })
		} else if t.terminal == tLiteral {
			nodes = append(nodes, &node{ value: t.num })
		} else if t.terminal == tOperator {
			n := 2
			if t.lexeme != "!" { n = 1 }
			if len(nodes) < n {
				log.Fatal("Invalid expression, not enough arguments")
			}
			root := &node{ nonterm: nExpr }
			for j := 0; j < n; j++ {
				root.addChild(nodes[len(nodes) - 1])
				nodes = nodes[:len(nodes) - 1]
			}
			subExpr = len(nodes)
			nodes = append(nodes, root)
			p.list = p.list[i:]
			i = 0
			continue
		} else {
			break
		}
		i++
	}
	if subExpr != 0 {
		if subExpr == -1 {
			return false
		} else if len(nodes) == 1 {
			root := &node{ nonterm: nExpr }
			root.addChild(nodes[0])
			nodes[0] = root
		} else {
			log.Fatal("Invalid expression, missing operator")
		}
	}
	p.cur = nodes[0]
	return true
}

func (p *parser) parseFunc_call() bool {
	if t, ok := p.getToken(0); !ok || t.lexeme != "(" { return false }
	p.list = p.list[1:]
	root := &node{ nonterm: nFunc_call }
	for p.parseExpr() {
		root.addChild(p.cur)
	}
	if len(root.child) == 0 || root.child[len(root.child) - 1].symbol == "" {
		log.Fatal("Invalid function call: no identifier")
	}
	if t, ok := p.getToken(0); !ok || t.lexeme != ")" {
		log.Fatal("Invalid function call")
	}
	p.list = p.list[1:]
	p.cur = root
	return true
}

