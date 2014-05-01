package main

import (
	"os"
	"fmt"
	"log"
)

func main() {
	tokens, err := lex(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	for t := range(tokens) {
		fmt.Println(t)
	}
}
