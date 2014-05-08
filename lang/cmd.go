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
	if t.lexeme == "" {
		fmt.Printf("lex %d: %d\n", t.terminal, t.num)
	} else {
		fmt.Printf("lex %d: %s\n", t.terminal, t.lexeme)
	}
}

func printChild(n *node, indent int) {
	tab := strings.Repeat("  ", indent)
	t := n.token
	if t.lexeme == "" {
		fmt.Printf("%s%d: %d\n", tab, n.nonterm, t.num)
	} else {
		fmt.Printf("%s%d: %s\n", tab, n.nonterm, t.lexeme)
	}
	for i := range(n.child) {
		printChild(n.child[i], indent+1)
	}
}
