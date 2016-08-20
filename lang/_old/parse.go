package lang

import (
	"fmt"
	"log"
	"strconv"
)

type nonterm int

const (
	nError nonterm = iota

	nFuncDef
	nBlock
	nIfStmt
	nAutoVarStmt
	nAliasStmt
	nAssignStmt
	nJumpStmt
	nReturnStmt

	nLabel
	nParam
	nType
	nExpr
	nFuncCall
)

type node struct {
	parent *node
	child  []*node
	nonterm
	*token
}

func (n *node) addChild(c *node) {
	c.parent = n
	n.child = append(n.child, c)
	return
}

type parser struct {
	t      chan token
	tokens []token
	tCur   int
	tree   *node
	sym    map[string]struct {
		size int
		*node
	}
	curNameSpace []string
}

func (p *parser) parseErr(t *token, err string) {
	if t != nil {
		log.Fatalf("Line %d, col %d: parsing error:\n\t%s", t.line, t.col, err)
	} else {
		log.Fatalf("Parsing error:\n\t%s", err)
	}
}

// Adds symbol to p.sym table
func (p *parser) addSymbol(n *node) {
	var sym string
	for _, s := range p.curNameSpace {
		sym = fmt.Sprintf("%s%s.", sym, s)
	}
	sym += n.lexeme
	if _, ok := p.sym[sym]; ok {
		p.parseErr(n.token, fmt.Sprintf("Redeclared symbol '%s'", sym))
	}
	s := p.sym[sym]
	s.node = n
	s.size = 8
	if n.parent != nil && n.parent.token != nil && n.parent.terminal == tLiteral {
		s.size, _ = strconv.Atoi(n.parent.lexeme)
	}
}

func parse(t chan token) *node {
	p := parser{t: t}
	p.tree = p.parseFile()
	if p.tree == nil {
		p.parseErr(nil, "Empty file")
	}
	return p.tree
}

func (p *parser) nextToken() bool {
	t, ok := <-p.t
	if !ok {
		return false
	}
	printToken(t)
	p.tokens = append(p.tokens, t)
	return true
}

func (p *parser) getToken(i int) (*token, bool) {
	if d := p.tCur + i - len(p.tokens); d >= 0 {
		for j := 0; j <= d; j++ {
			if !p.nextToken() {
				return nil, false
			}
		}
	}
	return &p.tokens[p.tCur+i], true
}

func (p *parser) parseFile() (n *node) {
	n = new(node)
	for {
		if t, ok := p.getToken(0); ok && t.terminal == tNewline {
			p.tCur++
			continue
		}
		c := p.parseStmt()
		if c == nil {
			break
		}
		n.addChild(c)
	}
	if len(n.child) == 0 {
		return nil
	}
	return
}

func (p *parser) parseStmt() (n *node) {
	stmt := func(f func() *node) bool {
		if s := f(); s != nil {
			if t, ok := p.getToken(0); ok && t.terminal != tNewline {
				return false
			}
			p.tCur++
			if n != nil {
				n.addChild(s)
			} else {
				n = s
			}
			return true
		}
		return false
	}
	tSaved := p.tCur
	if stmt(p.parseAliasStmt) {
		return
	}
	n = p.parseLabel()
	if n != nil {
		p.curNameSpace = append(p.curNameSpace, n.lexeme)
		if stmt(p.parseFuncDef) {
			return
		}
		p.curNameSpace = p.curNameSpace[:len(p.curNameSpace)-1]
		if stmt(p.parseBlock) || stmt(p.parseIfStmt) {
			return
		}
	}
	if stmt(p.parseAutoVarStmt) || stmt(p.parseAliasStmt) || stmt(p.parseAssignStmt) || stmt(p.parseJumpStmt) || stmt(p.parseReturnStmt) || stmt(p.parseParam) || stmt(p.parseExpr) {
		return
	}
	p.tCur = tSaved
	return nil
}

func (p *parser) parseFuncDef() (n *node) {
	t, ok := p.getToken(0)
	if !ok || t.terminal != tFunc {
		return
	}
	p.tCur++
	n = &node{nonterm: nFuncDef}
	if c := p.parseParam(); c != nil {
		n.addChild(c)
	}
	if t, ok := p.getToken(0); ok && t.terminal == tMap {
		p.tCur++
		if c := p.parseParam(); c != nil {
			n.addChild(c)
		} else {
			p.parseErr(t, "Invalid function definition")
		}
	}
	t, ok = p.getToken(0)
	// Block node
	if ok && t.terminal == tNewline {
		p.tCur++
		if c := p.parseBlock(); c != nil {
			n.addChild(c)
			return
		}
	}
	p.parseErr(t, "Invalid function definition")
	return nil
}

func (p *parser) parseBlock() (n *node) {
	t, ok := p.getToken(0)
	if !ok || t.terminal != tIndent {
		return
	}
	p.tCur++
	if n = p.parseFile(); n == nil {
		p.parseErr(t, "Invalid block statement")
	}
	n.nonterm = nBlock
	if t, ok := p.getToken(0); ok && t.terminal != tDedent {
		p.parseErr(t, "Invalid block statement")
	}
	p.tCur++
	return
}

func (p *parser) parseIfStmt() (n *node) {
	t, ok := p.getToken(0)
	if !ok || t.terminal != tIf {
		return
	}
	p.tCur++
	n = &node{nonterm: nIfStmt}
	if c := p.parseExpr(); c != nil {
		n.addChild(c)
	} else {
		p.parseErr(t, "Invalid if-statement")
	}
	// Block node
	if c := p.parseBlock(); c != nil {
		n.addChild(c)
	} else {
		p.parseErr(t, "Invalid if-statement")
	}
	return
}

func (p *parser) parseAutoVarStmt() (n *node) {
	tSaved := p.tCur
	n = &node{nonterm: nAutoVarStmt}
	if c := p.parseParam(); c != nil {
		n.addChild(c)
	} else {
		return nil
	}
	if t, ok := p.getToken(0); !ok || t.terminal != tAutoVar {
		p.tCur = tSaved
		return nil
	}
	p.tCur++
	if c := p.parseExpr(); c != nil {
		n.addChild(c)
	} else {
		p.parseErr(&p.tokens[p.tCur], "Invalid automatic variable definition")
	}
	return
}

func (p *parser) parseAliasStmt() (n *node) {
	tSaved := p.tCur
	n = &node{nonterm: nAliasStmt}
	if c := p.parseParam(); c != nil {
		n.addChild(c)
	} else {
		return nil
	}
	if t, ok := p.getToken(0); !ok || t.terminal != tAlias {
		p.tCur = tSaved
		return nil
	}
	if c := p.parseExpr(); c != nil {
		n.addChild(c)
	} else {
		p.tCur = tSaved
		return nil
	}
	p.tCur++
	return
}

func (p *parser) parseAssignStmt() (n *node) {
	numArgs := 0
	parseSingle := func(c *node) bool {
		t, ok := p.getToken(0)
		if !ok || t.terminal != tIdentifier {
			return false
		}
		c.addChild(&node{token: t})
		p.tCur++
		numArgs++
		return true
	}
	tSaved := p.tCur
	c := new(node)
	n = &node{nonterm: nAssignStmt}
	n.addChild(c)
	if !parseSingle(c) {
		return nil
	}
	for {
		if t, ok := p.getToken(0); ok && t.terminal == tComma {
			p.tCur++
			if !parseSingle(c) {
				p.parseErr(t, "Invalid assignment statement")
			}
		} else {
			break
		}
	}
	t, ok := p.getToken(0)
	if !ok {
		if numArgs == 1 {
			p.tCur = tSaved
			return nil
		}
		p.parseErr(&p.tokens[p.tCur], "Invalid parameter list")
	} else if t.terminal != tAssign {
		p.tCur = tSaved
		return nil
	}
	p.tCur++
	if c := p.parseExpr(); c != nil {
		n.addChild(c)
	} else {
		p.parseErr(t, "Invalid assignment statement")
	}
	return
}

func (p *parser) parseJumpStmt() (n *node) {
	if t, ok := p.getToken(0); !ok || t.terminal != tJump {
		return
	}
	p.tCur++
	n = &node{nonterm: nJumpStmt}
	if c := p.parseExpr(); c != nil {
		n.addChild(c)
	} else {
		p.parseErr(&p.tokens[p.tCur], "Invalid jump statement")
	}
	return
}

func (p *parser) parseReturnStmt() (n *node) {
	if t, ok := p.getToken(0); !ok || t.terminal != tReturn {
		return
	}
	p.tCur++
	return &node{nonterm: nReturnStmt}
}

func (p *parser) parseLabel() (n *node) {
	if t, ok := p.getToken(0); ok && t.terminal == tIdentifier {
		if u, ok := p.getToken(1); ok && u.terminal == tAlias {
			n = &node{
				nonterm: nLabel,
				token:   t,
			}
			p.tCur += 2
			p.addSymbol(n)
		}
	}
	return
}

func (p *parser) parseParam() (n *node) {
	n = &node{nonterm: nParam}
	parseSingle := func() bool {
		t, ok := p.getToken(0)
		if !ok || (t.terminal != tByte && t.terminal != tBlock && t.terminal != tIdentifier) {
			return false
		}
		child := n
		if t.terminal == tByte {
			child.addChild(&node{nonterm: nType, token: t})
			child = child.child[0]
			p.tCur++
		} else if t.terminal == tBlock {
			if u, ok := p.getToken(1); ok && u.terminal == tLiteral {
				child = &node{nonterm: nType, token: t}
				child = child.child[0]
				child.addChild(&node{nonterm: nType, token: u})
				child = child.child[0]
				p.tCur += 2
			} else {
				p.parseErr(t, "Invalid block-typed parameter")
			}
		}
		t, ok = p.getToken(0)
		if !ok || t.terminal != tIdentifier {
			p.parseErr(t, "Invalid parameter")
		}
		child.addChild(&node{token: t})
		p.addSymbol(child.child[0])
		p.tCur++
		return true
	}
	if !parseSingle() {
		return nil
	}
	for {
		if t, ok := p.getToken(0); ok && t.terminal == tComma {
			p.tCur++
			if !parseSingle() {
				p.parseErr(t, "Invalid parameter list")
			}
		} else {
			break
		}
	}
	return
}

func (p *parser) parseExpr() (n *node) {
	tSaved := p.tCur
	nodes := make([]*node, 0)
loop:
	for {
		t, ok := p.getToken(0)
		if !ok {
			break
		}
		if t.terminal == tLeftParen {
			nodes = append(nodes, p.parseFuncCall())
			p.tCur++
			continue
		}
		p.tCur++
		switch arity := 2; t.terminal {
		case tIdentifier:
			fallthrough
		case tLiteral:
			nodes = append(nodes, &node{token: t})
		case tNot:
			arity = 1
			fallthrough
		case tAdd:
			fallthrough
		case tSub:
			fallthrough
		case tMult:
			fallthrough
		case tDiv:
			fallthrough
		case tExp:
			fallthrough
		case tMod:
			fallthrough
		case tAnd:
			fallthrough
		case tOr:
			fallthrough
		case tXor:
			fallthrough
		case tShiftL:
			fallthrough
		case tShiftR:
			c := &node{nonterm: nExpr, token: t}
			for i := len(nodes) - arity - 1; i < len(nodes); i++ {
				c.addChild(nodes[i])
			}
			nodes[len(nodes)-arity-1] = c
			nodes = nodes[:len(nodes)-arity-1]
		default:
			break loop
		}
	}
	if len(nodes) == 0 {
		p.tCur = tSaved
		return nil
	} else if len(nodes) > 1 {
		p.parseErr(&p.tokens[p.tCur], "Invalid expression")
	}
	return nodes[0]
}

func (p *parser) parseFuncCall() (n *node) {
	if t, ok := p.getToken(0); !ok || t.terminal != tLeftParen {
		return
	}
	p.tCur++
	n = &node{nonterm: nFuncCall}
	for {
		if c := p.parseExpr(); c != nil {
			n.addChild(c)
		} else {
			break
		}
	}
	if len(n.child) == 0 || n.child[len(n.child)-1].child[0].terminal != tIdentifier {
		p.parseErr(&p.tokens[p.tCur-1], "Invalid function call: no identifier")
	}
	if t, ok := p.getToken(0); !ok || t.terminal != tRightParen {
		p.parseErr(t, "Invalid function call")
	}
	p.tCur++
	return
}
