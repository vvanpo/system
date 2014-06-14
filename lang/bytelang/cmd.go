// Bytelang virtual-machine and (de)compiler
package main

import (
	_ "os"
)

func main() {
	b := new(bytelang)
	b.statement = append(b.statement, bReturn)
}
