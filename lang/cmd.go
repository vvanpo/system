package main

import (
	"os"
	"fmt"
	//"log"
)

func main() {
	tokens := lex(os.Stdin)
	for t := range(tokens) {
		fmt.Printf("%#v\n", t)
	}
}
