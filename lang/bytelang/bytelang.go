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

// Literals are encoded as sequences of bytes, until the minimum word-length (bytelang.wordLength) can be established
// The byte slice uses big-endian ordering, which is also the representation used in the compiled bytecode file
type literal []byte

type symbol struct {
	*identifier
	address uint    // Address is an offset to the base namespace address, as pointed to by symbol.parent (0 if parent = nil)
	parent  *symbol // Parent namespace
	object			// Syntax object
}

type object interface {
	action()
}

func (s *symbol) action() {
	// TODO: add itself to symbol list
}

type variable struct {
	*symbol
	length // Length in bytes
}

func (v *variable) action() {
	// TODO: initialize as zero
}

type function struct {
	*symbol
	local []*variable // Local variables
	ret   []*variable // Return variables
	stmt  []*object // Statement list
}

func (f *function) action() {
	// TODO: 
}

type expression interface {
	evaluate()
}

type ifStmt struct {
	expr	expression
	stmt	[]statement
}

type assignmentStmt struct {
	
}

/*

func newBytecodeFile() (b *bytelang) {
	b = new(bytelang)
	b.identifier = []identifier{"_", "_sp", "_fp", "_ip", "_text", "_data"}
	for _, i := range b.identifier {
		b.addSymbol(i, 0, nil)
	}
	return
}

func (b *bytelang) addIdentifier(ident identifier) *identifier {
	for _, i := range b.identifier {
		if i == ident {
			return &i
		}
	}
	b.identifier = append(b.identifier, ident)
	return &b.identifier[len(b.identifier)-1]
}

func (b *bytelang) addSymbol(ident identifier, address uint, parent *symbol) (s *symbol) {
	s = new(symbol)
	s.identifier = b.addIdentifier(ident)
	s.address = address
	s.parent = parent
	b.symbol = append(b.symbol, *s)
	return
}

// Assumes big-endian ordering of words (most-significant word passed in as first argument)
func (b *bytelang) addWordLiteral(l ...uint) {
	lit := convertBigEndian(l...)
	for _, v := range b.literal {
		if bytes.Equal(lit, v) {
			return
		}
	}
	b.literal = append(b.literal, lit)
}

// Assumes big-endian ordering of bytes
func (b *bytelang) addByteLiteral(l ...byte) {
	for _, v := range b.literal {
		if bytes.Equal(l, v) {
			return
		}
	}
	b.literal = append(b.literal, l)
}

func (b *bytelang) addWord(w uint) {
	w = convertBigEndian(w)
	for uint(len(w)) < b.wordLength/8 {
		w = append([]byte{0}, w...)
	}
	b.code = append(b.code, w...)
}

func (b *bytelang) addSymbolDef(s string) {
	i := uint(len(b.symbol))
	b.currentNamespace = b.addSymbol(identifier(s), uint(len(b.code)), b.currentNamespace)
	b.code = append(b.code, bSymbolDef)
	b.addWord(i)
}

func convertBigEndian(n ...uint) []byte {
	num := []byte{}
	for _, b := range n {
		tmp := []byte{}
		for b != 0 {
			tmp = append([]byte{byte(0xff | b)}, tmp...)
			b >>= 8
		}
		num = append(num, tmp...)
	}
	return num
}
*/
