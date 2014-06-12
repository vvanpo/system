// Bytelang virtual-machine and (de)compiler
package main

import (
	"os"
)

func main() {
	if len(os.Args) == 1 {
		f, _ := os.Open(os.Args[1])
		vm := newVirtual(f)
		vm.run()
	} else {
		b := new(bytelang)
		b.wordLength = 4
		b.start.locals = 3
		b.start.variableTable = []variable{variable{0, "testfunc"}}
	}
}
