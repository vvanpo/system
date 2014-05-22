package main

type bytecode byte

// Bytecodes are ordered by expected rate of incidence to aid compression
const (
	bError     bytecode = iota
	bSymbolDef          // Namespace decided by nesting within bFunction declarations or bBlock* structures
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

//	Special values are special references via the symbol table, starting from the beginning of the table
//	bDiscard				//	Assignment results in nothing, identifier "_"
//	bStackPointer			//	Top of stack special register, identifier "_sp"
//	bFramePointer			//	Frame pointer special register, identifier "_fp"
//	bInstructionPointer		//	Current instruction special register, identifer "_ip"
//	bTextPointer			//	Beginning of text segment special variable "_text"
//	bDataPointer			//	Beginning of data/heap segment special variable "_data"
