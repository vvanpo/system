package main

const (
	bError byte = iota
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
	bVariableRef   // Value reference
	bFunctionCall  // First expression must be a variable or address
	bReferenceOp   // Address reference for a given variable
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

// The special identifiers belong to the top namespace
//	bDiscard				//	Assignment results in no-op, identifier "_"
//	bStackPointer			//	Top of stack special register, identifier "_sp"
//	bFramePointer			//	Frame pointer special register, identifier "_fp"
//	bInstructionPointer		//	Current instruction special register, identifer "_ip"
//	bTextPointer			//	Beginning of text segment special variable "_text"
//	bDataPointer			//	Beginning of data/heap segment special variable "_data"
const specialIdentifiers = "_\n_sp\n_fp\n_ip\n_text\n_data\n"

// Representation of a bytelang file
type bytelang struct {
	wordLength int          // Bytes per word
	identifier []identifier // Identifier list
	variable   []variable   // Variable list, indexes into identifier list
	imported   []*variable  // Imported variables, indexes into variable list
	start      function     // Top-level scope, program exit on return
	literal    []literal    // List of literals
}

type identifier string

type variable struct {
	*identifier
	parent    *variable // Parent namespace
	scope     *function
	refLength int  // Reference granularity in bytes, e.g. for bWord refLength = p.wordLen
	length    uint // Length in bytes, same as refLength if bByte or bWord
}

type function struct {
	bind  *variable   // Binding variable
	param []*variable // Parameters
	ret   []*variable // Return variables
	local []*variable // Local variables
	stmt  []statement // Statement list
}

type statement interface {
	action()
}

type ifStmt struct {
	expr []*expression // Condition
	stmt []statement   // Statement body
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

// Literals, as with all variable values, are encoded as sequences of bytes
// The byte slice uses big-endian ordering, which is also the representation used in the compiled bytecode file
type literal []byte

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
