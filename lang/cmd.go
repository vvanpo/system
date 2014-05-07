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

func printChild(n *node, indent int) {
	t := strings.Repeat("  ", indent)
	var val string
	if n.value != 0 {
		val = string(n.value)
	} else {
		val = n.symbol
	}
	fmt.Printf("%s%d: %s", t, n.nonterm, val)
	for i := range(n.child) {
		printChild(n.child[i], indent+1)
	}
}
