package main

// Bytecode markers
const (
	bError byte = iota
	bByte       // Byte-wide variable type
	bWord       // Word-wide variable type
	bBlock      // Compile-time-defined-length variable type
	// Statements:
	bVariable     // Variable allocation
	bFunctionCall // Function call statement
	bIf           // If statement
	bAssignment   // Assignment statement
	bReturn       // Return statement
	// Expressions:
	bLiteral   // Literal value
	bReference // Variable reference
	bFunction  // Function definition
	// Unary operators:
	bAddressReferenceOp // Variable address reference operator
	bDereferenceOp      // Variable dereference operator
	bNotOp
	// Binary operators:
	bAddOp
	bSubtractOp
	bMultiplyOp
	bDivideFloorOp
	bExponentOp
	bModuloOp
	bAndOp
	bOrOp
	bXorOp
	bRotateLeftOp
	bRotateRightOp
)

// The special identifiers are unique to the special variables, which have
// different behaviour than regular variables:
//		'_'		Assignment is discarded.
//		'_sp'	Stack pointer: the first operation in every program must be
//				setting the stack pointer.  _sp is updated with each
//				allocation and deallocation of the stack.
//		'_fp'	Frame pointer: like the stack pointer, _fp is automatically
//				updated on function calls.
//		'_ip'	Instruction pointer: can be used to implement jumps by
//				assigning address to _ip
const specialIdentifiers = "_\n_sp\n_fp\n_ip\n"

// Representation of a bytelang file
type bytelang struct {
	wordLength int          // Bytes per word
	identifier []identifier // Identifier list
	statement  []statement  // List of statements
}

func (b *bytelang) putWord(w uint) string {
	s := make([]byte, b.wordLength)
	for i := b.wordLength; i >= 0; i-- {
		s[i] = byte(0xff | w)
		w >>= 8
	}
	return string(s)
}

func (b *bytelang) identifierIndex(i *identifier) uint {
	for j := range b.identifier {
		if &b.identifier[j] == i {
			return uint(j + 1)
		}
	}
	return 0
}

func (b *bytelang) statementIndex(l statement) uint {
	for i := range b.statement {
		if b.statement[i] == l {
			return uint(i + 1)
		}
	}
	return 0
}

type identifier string

type bytecode interface {
	bytecode() string // All bytecode representations start with a bytecode marker (see above const list)
}

type statement interface {
	bytecode
}

type variable struct {
	*bytelang
	*identifier
}

type variableWord variable

func (v *variableWord) bytecode() (b string) {
	b = string(bVariable)
	b += v.putWord(v.identifierIndex(v.identifier))
	b += string(bWord)
	return
}

type variableByte variable

func (v *variableByte) bytecode() (b string) {
	b = string(bVariable)
	b += v.putWord(v.identifierIndex(v.identifier))
	b += string(bByte)
	return
}

type variableBlock struct {
	variable
	length uint // Length in bytes
	member []struct {
		offset uint
	}
}

func (v *variableBlock) bytecode() (b string) {
	b = string(bVariable)
	b += v.putWord(v.identifierIndex(v.identifier))
	b += string(bBlock)
	b += v.putWord(v.length) + v.putWord(uint(len(v.member)))
	for _, m := range v.member {
		b += v.putWord(m.offset)
	}
	return
}

type functionCall struct {
	*bytelang
	callee   expression
	argument []expression
	receiver []statement
}

func (f *functionCall) bytecode() (b string) {
	b = string(bFunctionCall) + f.callee.bytecode()
	for _, a := range f.argument {
		b += a.bytecode()
	}
	for _, r := range f.receiver {
		b += f.putWord(f.statementIndex(r))
	}
	return
}

type ifStmt struct {
	*bytelang
	condition expression
	statement []statement
}

func (i *ifStmt) bytecode() (b string) {
	b = string(bIf)
	b += i.condition.bytecode()
	for _, s := range i.statement {
		b += i.putWord(i.statementIndex(s))
	}
	return
}

type assignment struct {
	*bytelang
	receiver statement
	expression
}

func (a *assignment) bytecode() (b string) {
	b = string(bAssignment)
	b += a.putWord(a.statementIndex(a.receiver))
	b += a.expression.bytecode()
	return
}

type returnStmt struct{}

func (r *returnStmt) bytecode() (b string) {
	b = string(bReturn)
	return
}

type expression interface {
	bytecode
	value() []byte
}

type literal struct {
	*bytelang
	bytes []byte
}

func (l *literal) bytecode() (b string) {
	b = string(bLiteral)
	b += l.putWord(uint(len(l.bytes)))
	b += string(l.bytes)
	return
}

func (l *literal) value() (b []byte) {
	b = l.bytes
	return
}

type reference struct {
	*bytelang
	statement
}

func (r *reference) bytecode() (b string) {
	b = string(bReference)
	b += r.putWord(r.statementIndex(r.statement))
	return
}

func (r *reference) value() (b []byte) {
	return
}

type function struct {
	*bytelang
	parameter []statement
	returnVal []statement
	statement []statement
}

func (f *function) bytecode() (b string) {
	b = string(bFunction)
	return
}

func (f *function) value() (b []byte) {
	return
}
