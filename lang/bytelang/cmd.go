// Bytelang virtual-machine and (de)compiler
package main

import (
	_ "os"
)

func (f *function) add(s ...statement) {
	*f = append(*f, s...)
}

func main() {
}
