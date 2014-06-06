package main

const (
	bError    byte = iota
	bBytes         // Bytes literal
	bFunction      // Function literal
	bByte          // Variable types
	bWord
	bBlock
	bVariable // Statement types
	bFunctionCall
	bIf
	bAssignment
	bReturn
	bLiteralPointer     // Pointer to literal
	bLiteralValue       // Value-copy of literal
	bReference          // Variable reference
	bAddressReferenceOp // Variable address reference operator
	bDereferenceOp      // Address dereference operator
	bNotOp
	bAddOp // Binary operators
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
	literal    []literal    // List of literals
}

type identifier string

type namespace struct {
	*variable
	*identifier // Index into identifier list
	member      []namespace
}

type literal interface {
	value() string
}

type function struct {
	parameter	[]variable
	returnVal	[]variable
	statement	[]statement
}

func (f *function) value() string {

}

type bytes string

func (b bytes) value() string {
	return b
}

type variable struct {
	length    uint      // Length in bytes
}

type statement interface {
	
}

type expression interface {

}

/*
type variable struct {
	*identifier
	scope     *function
	base      *variable // Base address
	// Automatic variables will use _fp
	// Address aliases (heap variables) will use _data
	// Function aliases will use _text
	offset int // Address offset from base
}

type function struct {
	bind  *variable   // Binding variable
	param []*variable // Parameters
	local []*variable // Local variables
	stmt  []statement // Statement list
}

type statement interface {
	action()
}

type variableDefStmt struct{}

func (v *variableDefStmt) action() {
}

type ifStmt struct {
	expr []expression // Condition
	stmt []statement  // Statement body
}

func (i *ifStmt) action() {
}

type assignmentStmt struct {
	assignee []*variable
	expr     expression
}

func (a *assignmentStmt) action() {
}

type jumpStmt struct {
	expr expression
}

func (j *jumpStmt) action() {
}

type returnStmt struct{}

func (r *returnStmt) action() {
}

type expression interface {
	evaluate() []byte
}

func (l literal) evaluate() []byte {
	return l
}

type reference struct {
	*variable
}

type functionCall *function

type operation struct {
	expr []*expression
}
*/
