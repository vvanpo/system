package main

// Bytecode markers
const (
	bError byte = iota
	bByte       // Byte-wide variable type
	bWord       // Word-wide variable type
	bBlock      // Compile-time-defined-length variabl type
	// Statements:
	bVariable     // Variable declaration
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
	namespace               // Namespace tree
	statement  []statement  // List of literals
}

func (b *bytelang) putWord(w uint) string {
	s := make([]byte, b.wordLength)
	for i := b.wordLength; i >= 0; i-- {
		s[i] = byte(0xff | w)
		w >>= 8
	}
	return string(s)
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

type namespace struct {
	*variable   // Reverse index into statement list
	*identifier // Index into identifier list
	member      []*namespace
	parent      *namespace
}

type bytecode interface {
	bytecode() (marker byte, bytecode string) // All bytecode representations start with a bytecode marker (see above const list)
}

type statement interface {
	bytecode
}

type variable struct {
	*bytelang
	typ   byte // bByte, bWord, or bBlock
	block struct {
		length uint // Length in bytes
		member []struct {
			offset uint
			*namespace
		}
	}
}

func (v *variable) bytecode() (marker byte, bytecode string) {
	marker = bVariable
	bytecode = string(v.typ)
	if v.typ == bBlock {
		bytecode += v.putWord(v.block.length) + v.putWord(uint(len(v.member)))
		for _, m := range v.block.member {
			bytecode += v.putWord(m.offset)
			bytecode += v.putWord(v.statementIndex(m.variable))
		}
	}
	return
}

type functionCall struct {
	*bytelang
	callee   expression
	argument []expression
	receiver []*namespace
}

func (f *functionCall) bytecode() (marker byte, bytecode string) {
	marker = bFunctionCall
	bytecode = f.callee.bytecode()
	for _, a := range f.argument {
		bytecode += a.bytecode()
	}
	for _, r := range f.receiver {
		bytecode += f.putWord(f.namespaceIndex(r))
	}
}

type expression interface {
	bytecode
	value() []byte
}

type function struct {
	parameter []*namespace
	returnVal []*namespace
	statement []*statement
}
