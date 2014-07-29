
from pyparsing import *

def parse(code):
	ParserElement.setDefaultWhitespaceChars(" \t\r")
	name = Word(alphas)
	number = Word(hexnums)
	offset = number
	uuid = number
	address = (name + Suppress(":") + number) | (oneOf("_ip _sp _fp") + Optional(offset))
	length = number | Literal("-")
	call = "call" + address
	ref = "ref" + address
	deref = "deref" + address + length
	literal = "literal" + number + length
	operation = oneOf("not and or xor shiftl lshiftr ashiftr add sub mult \
			floordiv exp mod") + length
	expression = call | ref | deref | literal | operation
	openfd = "open" + (name ^ (Optional(name) + uuid))
	close = "close" + (name | uuid)
	push = "push" + expression
	pop = "pop" + Optional(number) + Optional(address + Optional(length + offset))
	copy = "copy" + expression + address
	ifzero = "ifzero" + expression + address
	statement = (openfd | close | push | pop | copy | ifzero) + lineEnd.suppress()
	statement.setWhitespaceChars(" \t\r\n")
	asmlang = OneOrMore(statement) + stringEnd

	instr = []
	def action(s, loc, toks):
		instr.append(toks.asList())
	statement.setParseAction(action)
	asmlang.parseString(code)
	return instr
