// Bytelang virtual-machine and (de)compiler
package main

import (
	"os"
)

func main() {
	f, _ := os.Open(os.Args[1])
	vm := newVirtual(f)
	vm.run()
}
