// Bytelang virtual-machine and (de)compiler
package main

import (
	_ "os"
)

func (f *function) add(s ...statement) {
	f.statement = append(f.statement, s...)
}

func main() {
	b := new(bytelang)
	b.add(allocate(8))
	b.add(assignment{
		stackPointer{0},
		instructionPointer{},
		8,
	})
	c := b.compile()
	print(c)
}
