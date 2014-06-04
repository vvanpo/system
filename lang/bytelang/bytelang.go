package main

const (
	bError byte = iota
	bVariableDef
	bWord
	bByte
	bBlockWord
	bBlockByte
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
	scope     *function
	refLength int       // Reference granularity in bytes, e.g. for bWord refLength = bytelang.wordLength
	length    uint      // Length in terms of refLength
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
