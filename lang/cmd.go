package main

import (
	"os"
)

func main() {
	tokens := lex(os.Stdin)
	parse(tokens)
}
