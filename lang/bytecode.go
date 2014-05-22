package main

import (
	"bytes"
)

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

type identifier string

type symbol struct {
	*identifier
	address uint
	parent  *symbol
}

// Literals are encoded as sequences of bytes, until the minimum word-length (bytecodeFile.wordLength) can be established
// The byte slice uses big-endian ordering, which is also the representation used in the compiled bytecode file
type literal []byte

type bytecodeFile struct {
	wordLength  uint         // Minimum word-length is decided on using the largest memory
	identifier  []identifier // reference made in the compiled code.  Literals that exceed
	symbol      []symbol     // this value are converted to block types
	imported    []*identifier
	literal     []literal
	textPointer uint // Beginning of text segment special variable "_text"
	dataPointer uint // Beginning of data/heap segment special variable "_data"
}

func newBytecodeFile() (b *bytecodeFile) {
	b = new(bytecodeFile)
}

// Assumes big-endian ordering of words (most-significant word passed in as first argument)
func (b *bytecodeFile) addWordLiteral(l ...uint) {
	lit := []byte{}
	for _, i := range l {
		tmp := []byte{}
		for i != 0 {
			tmp = append([]byte{byte(0xff | i)}, tmp...)
			i >>= 8
		}
		lit = append(lit, tmp...)
	}
	for _, v := range b.literal {
		if bytes.Equal(lit, v) {
			return
		}
	}
	b.literal = append(b.literal, lit)
}

// Assumes big-endian ordering of bytes
func (b *bytecodeFile) addByteLiteral(l ...byte) {
	for _, v := range b.literal {
		if bytes.Equal(l, v) {
			return
		}
	}
	b.literal = append(b.literal, l)
}
