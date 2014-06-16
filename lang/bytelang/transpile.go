package bytelang

// Transpile a bytelang structure into C89 source code
func (b *Bytelang) Transpile() (s string) {
	s = b.function.transpile()
	return
}

func (f function) transpile() (s string) {
	return
}
