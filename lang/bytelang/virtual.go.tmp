package main

import (
	"io"
)

type virtual struct {
	*parser
}

// New virtual machine
func newVirtual(r io.Reader) (vm *virtual) {
	vm.parser = newParser(r)
	return
}

func (vm *virtual) run() {
	vm.parseBytelang()
}
