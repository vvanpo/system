package main

// Bytecode markers
const (
	// Statements:
	bAssignment   byte = iota // Assignment statement
	bFunctionCall             // Function call statement
	bThread
	bIf
	bReturn
	// Expressions:
	bFunction  // Function definition
	bReference // Variable reference
	bLiteral   // Literal value
	// Unary operators:
	bDereferenceOp // Variable dereference operator
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

// Representation of a bytelang file
type bytelang struct {
	wordLength int // Bytes per word
	global     function
}

type function struct {
	parent           *function
	callerAllocation uint // Allocation size in words
	allocation       uint
	returns          []variable
	params           []variable
	locals           []variable
	statement        []interface{}
}

type variable struct {
	address    uint
	identifier string
}

type assignment struct {
	address    uint
	expression interface{}
}

type functionCall struct {
	address    uint
	allocation uint        // Allocation size in words
	returns    []*variable // Index into owning function's variable table
	args       []*variable
}

type thread functionCall

type ifStmt struct {
	condition interface{}
	statement []interface{}
}
