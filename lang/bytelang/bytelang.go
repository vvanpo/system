package main

// Bytecode markers
const (
	// Statements:
	bFunction byte = iota
	bAllocate
	bAssignment
	bThread
	bIf
	bReturn
	// Expressions:
	bFunctionCall
	bReference
	bDereference
	bLiteral
	bNot
	bAnd
	bOr
	bXor
	bShiftL
	bLShiftR
	bAShiftR
	bAdd
	bSubtract
	bMultiply
	bDivideFloor
	bExponent
	bModulo
)

// Representation of a bytelang file
type bytelang struct {
	wordLength int // Bytes per word
	global     function
}

type function struct {
	parent           *function
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
