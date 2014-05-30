package main

import (
	"io"
)

type virtual struct {
	*parser
}

// New virtual machine
func newVirtual(r io.Reader) (vm *virtual) {
	vm = new(virtual)
	vm.parser = newParser(r)
}

func (vm *virtual) run() {

}
