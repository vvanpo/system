package main

const (
	bError     byte = iota
	bSymbolDef      // Namespace decided by nesting within bFunction declarations or bBlock* structures
	bWord
	bByte
	bBlockWord
	bBlockByte
	bAutomatic // Allocate automatic variable on the stack
	bAddress   // Index into specific address value
	bOffset    // Must be an offset into a bBlock* structure.  Variable numbers start with parameters, then return values, then local variables
	bFunction
	bIf
	bAssignment
	bJump
	bReturn
	bExpression
	bLiteral
	bSymbolRef
	bFunctionCall // First expression must be a symbol or address
	bOperation
	bReferenceOp   // Address reference for a given symbol
	bDereferenceOp // Pointer dereference
	bAddOp
	bSubtractOp
	bMultiplyOp
	bDivideOp
	bExponentOp
	bModuloOp
	bAndOp
	bOrOp
	bXorOp
	bNotOp
	bRotateLeftOp
	bRotateRightOp
)

// The special identifiers belong to the current namespace, and have an address offset of 0
//	bDiscard				//	Assignment results in nothing, identifier "_"
//	bStackPointer			//	Top of stack special register, identifier "_sp"
//	bFramePointer			//	Frame pointer special register, identifier "_fp"
//	bInstructionPointer		//	Current instruction special register, identifer "_ip"
//	bTextPointer			//	Beginning of text segment special variable "_text"
//	bDataPointer			//	Beginning of data/heap segment special variable "_data"
const specialIdentifiers = "_\n_sp\n_fp\n_ip\n_text\n_data\n"

type identifier string

// Syntax object (function, statement, expression, etc.)
type object interface {
	action()
	up() *object // Return parent object
}

type symbol struct {
	*identifier
	address uint    // Address is an offset to the base namespace address, as pointed to by symbol.parent (0 if parent = nil)
	parent  *symbol // Parent namespace
}

type variable struct {
	*symbol
	scope  *function
	length int // Length in bytes
}

type function struct {
	name  *variable   // Binding variable (optional)
	param []*variable // Parameters
	ret   []*variable // Return variables
	stmt  []*object   // Statement list
}

type ifStmt struct {
	expr []*expression
	stmt []*object
}

type assignmentStmt struct {
	assignee []*variable
	expr     *expression
}

type jumpStmt struct {
	expr *expression
}

type returnStmt struct{}

type expression interface {
	evaluate() []byte
}

// Literals are encoded as sequences of bytes, until the minimum word-length (bytelang.wordLength) can be established
// The byte slice uses big-endian ordering, which is also the representation used in the compiled bytecode file
type literal []byte

func (l literal) evaluate() []byte {
	return l
}

type reference struct {
	*variable
}

func (r *reference) evaluate() []byte {
}

type functionCall *function

type operation struct{}
