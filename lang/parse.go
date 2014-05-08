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

type parser struct {
	tokens chan token
	list   []token
	tIndex int
	tree   *node
	cur    *node
	sym		map[string]struct{
		size	int
		*node
	}
}

type node struct {
	parent *node
	child  []*node
	nonterm
	token
}

func (p *parser) parseErr(t token, err string) {
	log.Fatalf("Parsing error at Line %d, Col %d:\n\t%s", t.line, t.col, err)
	return
}

func (p *parser) addSymbol(n *node) {
	sym := n.token.lexeme
	//sym = fmt.Sprintf("%s.%s", , sym)
	if _, ok := p.sym[sym] {
		p.parseErr(n.token, fmt.Sprintf("Redeclared symbol '%s'", sym))
	}
	size := 8
	p.sym[sym] = struct{ size, n }
}

func parse(tokens chan token) *node {
	p := parser{
		tokens: tokens,
		tree:   new(node),
	}
	p.tree.parent = p.tree
	p.cur = p.tree
	p.parseFile()
	return p.tree
}

func (p *parser) nextToken() bool {
	t, ok := <-p.tokens
	if !ok {
		return false
	}
	printToken(t)
	p.list = append(p.list, t)
	return true
}

func (p *parser) getToken(i int) (token, bool) {
	if d := p.tIndex + i - len(p.list); d >= 0 {
		for j := 0; j <= d; j++ {
			if !p.nextToken() {
				return token{}, false
			}
		}
	}
	return p.list[p.tIndex+i], true
}

func (n *node) addChild(c *node) {
	c.parent = n
	n.child = append(n.child, c)
	return
}

func (p *parser) parseFile() bool {
	root := p.cur
	if !p.parseStmt() {
		return false
	}
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
	p.cur = root
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
			param := &node{nonterm: nParam}
			param.addChild(label)
			root.addChild(param)
			// Add expression to label_stmt
			root.addChild(p.cur)
			p.cur = root
			return true
		}
	}
	if p.parseAssign_stmt() || p.parseLabel_stmt() || p.parseReassign_stmt() || p.parseJump_stmt() || p.parseReturn_stmt() {
		return true
	}
	// If the statement was invalid, check if a label was parsed
	if len(root.child) > 0 {
		p.parseErr(root.child[0].token, "Invalid statement")
	}
	p.cur = saved
	return false
}

func (p *parser) parseLabel() bool {
	if t, ok := p.getToken(0); ok && t.terminal == tIdentifier {
		if u, ok := p.getToken(1); ok && u.lexeme == ":" {
			p.cur = &node{
				parent:  p.cur,
				nonterm: nLabel,
				token:   t,
			}
			p.tIndex += 2
			p.addSymbol(t.lexeme, p.cur)
			return true
		}
	}
	return false
}

func (p *parser) parseFunc_def() bool {
	root := p.cur
	if t, ok := p.getToken(0); !ok || t.lexeme != "func" {
		return false
	}
	p.tIndex++
	for p.parseParam() {
		root.addChild(p.cur)
	}
	if t, ok := p.getToken(0); ok && t.lexeme == "->" {
		p.tIndex++
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
	t, ok := p.getToken(0)
	if !ok || t.terminal != tIndent {
		return false
	}
	p.tIndex++
	p.cur.nonterm = nBlock
	if !p.parseFile() {
		p.parseErr(t, "Invalid block statement")
	}
	if t, ok := p.getToken(0); ok && t.terminal != tDedent {
		p.parseErr(t, "Invalid block statement")
	}
	p.tIndex++
	return true
}

func (p *parser) parseIf_stmt() bool {
	if t, ok := p.getToken(0); !ok || t.lexeme != "if" {
		return false
	}
	p.tIndex++
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
		if !ok || (t.terminal != tReserved && t.terminal != tIdentifier) {
			return false
		}
		if t.terminal == tReserved {
			if t.lexeme == "byte" {
				root.addChild(&node{nonterm: nType, token: t})
				p.tIndex++
			} else if t.lexeme == "block" {
				if u, ok := p.getToken(1); ok && u.terminal == tLiteral {
					root.addChild(&node{nonterm: nType, token: t})
					root.addChild(&node{nonterm: nType, token: u})
					p.tIndex += 2
				} else {
					log.Fatal("Invalid block-typed parameter")
				}
			}
			t, ok = p.getToken(0)
			if !ok || t.terminal != tIdentifier {
				log.Fatal("Invalid parameter")
			}
		}
		root.addChild(&node{token: t})
		p.tIndex++
		return true
	}
	if !parseSingle() {
		return false
	}
	for {
		if t, ok := p.getToken(0); ok && t.lexeme == "," {
			p.tIndex++
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
	tIndex := p.tIndex
	defer func() { p.cur = root }()
	if !p.parseParam() {
		return false
	}
	param := p.cur
	if t, ok := p.getToken(0); !ok || t.lexeme != ":=" {
		p.tIndex = tIndex
		return false
	}
	p.tIndex++
	if !p.parseExpr() {
		p.tIndex = tIndex
		return false
	}
	root.addChild(param)
	root.addChild(p.cur)
	root.nonterm = nAssign_stmt
	return true
}

func (p *parser) parseLabel_stmt() bool {
	root := p.cur
	tIndex := p.tIndex
	defer func() { p.cur = root }()
	if !p.parseParam() {
		return false
	}
	param := p.cur
	if t, ok := p.getToken(0); !ok || t.lexeme != ":" {
		p.tIndex = tIndex
		return false
	}
	if !p.parseExpr() {
		p.tIndex = tIndex
		return false
	}
	p.tIndex++
	root.addChild(param)
	root.addChild(p.cur)
	root.nonterm = nLabel_stmt
	return true
}

func (p *parser) parseReassign_stmt() bool {
	root := p.cur
	list := make([]*node, 0)
	tIndex := p.tIndex
	parseSingle := func() bool {
		t, ok := p.getToken(0)
		if !ok || t.terminal != tIdentifier {
			return false
		}
		list = append(list, &node{token: t})
		p.tIndex++
		return true
	}
	if !parseSingle() {
		return false
	}
	for {
		if t, ok := p.getToken(0); ok && t.lexeme == "," {
			p.tIndex++
			if !parseSingle() {
				p.tIndex = tIndex
				return false
			}
		} else {
			break
		}
	}
	t, ok := p.getToken(0)
	p.tIndex++
	if !ok {
		if len(list) == 1 {
			p.tIndex = tIndex
			return false
		}
		log.Fatal("Invalid parameter list")
	} else if t.lexeme != "=" {
		p.tIndex = tIndex
		return false
	}
	if !p.parseExpr() {
		log.Fatal("Invalid assignment statement")
	}
	for _, n := range list {
		root.addChild(n)
	}
	root.addChild(p.cur)
	root.nonterm = nReassign_stmt
	return true
}

func (p *parser) parseJump_stmt() bool {
	if t, ok := p.getToken(0); !ok || t.lexeme != "jump" {
		return false
	}
	p.tIndex++
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
	if t, ok := p.getToken(0); !ok || t.lexeme != "return" {
		return false
	}
	p.tIndex++
	p.cur.nonterm = nReturn_stmt
	return true
}

func (p *parser) parseExpr() bool {
	nodes := make([]*node, 0)
	tIndex := p.tIndex
	subExpr := -1
	funcCall := -1
	for {
		t, ok := p.getToken(0)
		if !ok {
			break
		}
		if t.lexeme == "(" {
			p.parseFunc_call()
			nodes = append(nodes, p.cur)
			funcCall++
			continue
		}
		p.tIndex++
		if t.terminal == tIdentifier {
			nodes = append(nodes, &node{token: t})
		} else if t.terminal == tLiteral {
			nodes = append(nodes, &node{token: t})
		} else if t.terminal == tOperator {
			n := 2
			if t.lexeme == "!" {
				n = 1
			}
			if len(nodes) < n {
				log.Fatal("Invalid expression, not enough arguments")
			}
			root := &node{
				nonterm: nExpr,
				token:   t,
			}
			for j := 0; j < n; j++ {
				root.addChild(nodes[len(nodes)-1])
				nodes = nodes[:len(nodes)-1]
			}
			subExpr = len(nodes)
			nodes = append(nodes, root)
			continue
		} else {
			p.tIndex--
			break
		}
	}
	if subExpr == -1 && funcCall <= 0 {
		if len(nodes) >= 1 {
			root := &node{nonterm: nExpr}
			root.addChild(nodes[0])
			p.tIndex -= len(nodes) - 1
			p.cur = root
			return true
		} else {
			p.tIndex = tIndex
			return false
		}
	} else if subExpr == 0 {
		p.tIndex -= len(nodes) - 1
		p.cur = nodes[0]
		return true
	}
	p.parseErr(p.list[tIndex], "Invalid expression")
	return false
}

func (p *parser) parseFunc_call() bool {
	if t, ok := p.getToken(0); !ok || t.lexeme != "(" {
		return false
	}
	p.tIndex++
	root := &node{nonterm: nFunc_call}
	for p.parseExpr() {
		root.addChild(p.cur)
	}
	if len(root.child) == 0 || root.child[len(root.child)-1].child[0].token.terminal != tIdentifier || len(root.child[len(root.child)-1].child) != 1 {
		log.Fatal("Invalid function call: no identifier")
	}
	if t, ok := p.getToken(0); !ok || t.lexeme != ")" {
		log.Fatal("Invalid function call")
	}
	p.tIndex++
	p.cur = root
	return true
}
