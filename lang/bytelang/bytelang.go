package main

// Bytecode markers
const (
	bAddress byte = iota
	// Globals:
	bStackPointer
	bFramePointer
	bInstructionPointer
	// Statements:
	bFunction
	bAllocate
	bDeallocate
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
	function
}

type statement interface {
	compile() string
}

type function []statement

type allocate uint

type deallocate uint

type assignment struct {
	address
	value  expression
	length uint
}

type thread functionCall

type ifStmt struct {
	condition expression
	statement []statement
}

type returnStmt struct{}

type expression interface {
	compile() string
}

type functionCall uint

type reference uint

type dereference struct {
	address
	length uint
}

type literal []uint

type operation struct {
	marker byte
	length uint
}

type address interface {
	compile() string
}

type stackPointer struct {
	offset uint
}

type framePointer struct {
	offset uint
}

type instructionPointer struct{}
