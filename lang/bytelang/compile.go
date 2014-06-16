package main

func putWord(word uint) string {
	s := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		s[i] = byte(0xff & word)
		word >>= 8
	}
	return string(s)
}

func (b *bytelang) compile() (s string) {
	s = b.function.compile()[1:] // Omit the marker bytecode for the global function
	return
}

func (f function) compile() (s string) {
	s = string(bFunction)
	s += putWord(uint(len(f)))
	for _, stmt := range f {
		s += stmt.compile()
	}
	return
}

func (a allocate) compile() (s string) {
	s = string(bAllocate)
	s += putWord(uint(a))
	return
}

func (d deallocate) compile() (s string) {
	s = string(bDeallocate)
	s += putWord(uint(d))
	return
}

func (a assignment) compile() (s string) {
	s = string(bAssignment)
	s += a.address.compile()
	s += a.value.compile()
	s += putWord(a.length)
	return
}

func (t thread) compile() (s string) {
	s = string(bThread)
	s += putWord(uint(t))
	return
}

func (i ifStmt) compile() (s string) {
	s = string(bIf)
	s += i.condition.compile()
	for _, stmt := range i.statement {
		s += stmt.compile()
	}
	return
}

func (r returnStmt) compile() (s string) {
	s = string(bReturn)
	return
}

func (f functionCall) compile() (s string) {
	s = string(bFunctionCall)
	s += putWord(uint(f))
	return
}

func (r reference) compile() (s string) {
	s = string(bReference)
	s += putWord(uint(r))
	return
}

func (d dereference) compile() (s string) {
	s = string(bDereference)
	s += d.address.compile()
	s += putWord(d.length)
	return
}

func (l literal) compile() (s string) {
	s = string(bLiteral)
	for _, w := range l {
		s += putWord(w)
	}
	return
}

func (o operation) compile() (s string) {
	s = string(o.marker)
	s += putWord(o.length)
	return
}

func (sp stackPointer) compile() (s string) {
	s = string(bStackPointer)
	s += putWord(sp.offset)
	return
}

func (f framePointer) compile() (s string) {
	s = string(bFramePointer)
	s += putWord(f.offset)
	return
}

func (i instructionPointer) compile() (s string) {
	s = string(bInstructionPointer)
	return
}
