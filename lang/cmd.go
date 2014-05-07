package main

import (
	"os"
	"strings"
	"fmt"
)

func main() {
	tokens := lex(os.Stdin)
	tree := parse(tokens)
	printChild(tree, 0)
}

func printToken(t token) {
	var val string
	if t.lexeme == "" {
		val = string(t.num)
	} else {
		val = t.lexeme
	}
	fmt.Printf("lex %d: %s\n", t.terminal, val)
}

func printChild(n *node, indent int) {
	t := strings.Repeat("  ", indent)
	var val string
	if n.token.lexeme == "" {
		val = string(n.token.num)
	} else {
		val = n.token.lexeme
	}
	fmt.Printf("%s%d: %s", t, n.nonterm, val)
	for i := range(n.child) {
		printChild(n.child[i], indent+1)
	}
}
