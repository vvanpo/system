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

type function struct {
	parent    *function
	statement []statement
}

type allocate uint

type assignment struct {
	address
	value  address
	length uint
}

type thread functionCall

type ifStmt struct {
	condition interface{}
	statement []statement
}

type returnStmt struct{}

type functionCall uint

type reference uint

type dereference struct {
	address
	length uint
}

type literal []uint

type notOp struct {
	address
}

type andOp binaryOp
type orOp binaryOp
type xorOp binaryOp
type shiftLOp binaryOp
type lShiftROp binaryOp
type aShiftROp binaryOp
type addOp binaryOp
type subtractOp binaryOp
type multiplyOp binaryOp
type divideFloorOp binaryOp
type exponentOp binaryOp
type moduloOp binaryOp

type binaryOp struct {
	operandOne address
	operandTwo address
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
