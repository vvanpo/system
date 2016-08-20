package lang

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	tokens := lex(os.Stdin)
	tree := parse(tokens)
	printTree(tree)
}

func printToken(t token) {
	fmt.Printf("lex %d: %s\n", t.terminal, t.lexeme)
}

func printTree(n *node) {
	var recurse func(n *node, indent int)
	recurse = func(n *node, indent int) {
		tab := strings.Repeat("  ", indent)
		val := ""
		if n.token != nil {
			val = n.lexeme
		}
		fmt.Printf("%s%d: %s\n", tab, n.nonterm, val)
		for i := range n.child {
			recurse(n.child[i], indent+1)
		}
	}
	recurse(n, 0)
}
