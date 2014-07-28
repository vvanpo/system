
from pyparsing import *

def parse(code):
	name = Word(alphas)
	number = Word(hexnums)
	offset = number
	uuid = number
	address = (name + ":" + number) | (oneOf("_ip _sp _fp") + Optional(offset))
	length = number | Literal("-")
	call = "call" + address
	ref = "ref" + address
	deref = "deref" + address + length
	literal = "literal" + number + length
	operation = oneOf("not and or xor shiftl lshiftr ashiftr add sub mult \
			floordiv exp mod") + length
	expression = call | ref | deref | literal | operation
	open = "open" + (name | (Optional(name) + uuid))
	close = "close" + (name | uuid)
	push = "push" + expression
	pop = "pop" + Optional(number) + Optional(address + Optional(length + offset))
	copy = "copy" + expression + address
	ifzero = "ifzero" + expression + address
	statement = (open | close | push | pop | copy | ifzero) + lineEnd.suppress()
	asmlang = OneOrMore(statement.setDebug()) + stringEnd

	instr = []
	def action(s, loc, toks):
		instr.append(toks.asList())
	expression.setParseAction(action)
	statement.setParseAction(action)
	asmlang.parseString(code)
	return instr
