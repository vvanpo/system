
from pyparsing import *

def parse(code):
	name = Word(alphas)
	number = Word(hexnums)
	address = name + ":" + number | oneOf("_ip _sp _fp") + Optional(number)
	length = number | Literal("-")
	call = "call" + address
	ref = "ref" + address + length
	deref = "deref" + address + length
	literal = "literal" + number
	operation = oneOf("not and or xor shiftl lshiftr ashiftr add sub mult \
			divfloor exp mod") + length
	expression = call | ref | deref | literal | operation
	segment = "segment" + name
	pop = "pop" + Optional(number) + Optional(address)
	ifstmt = "if" + address
	statement = segment | pop | ifstmt
	asmlang = OneOrMore(statement | expression) + stringEnd

	instr = []
	def action(s, loc, toks):
		instr.append(toks.asList())
	expression.setParseAction(action)
	statement.setParseAction(action)
	asmlang.parseString(code)
	return instr
