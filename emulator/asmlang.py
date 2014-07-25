
from pyparsing import *

def parse(code):
	name = Word(alphas)
	number = Word(nums)
	address = name + ":" + number |
				(("_ip" | "_sp" | "_fp") + number)
	call = "call" + address
	ref = "ref" + address + number
	deref = "deref" + address + number
	literal = "literal" + number
	operation = ("not" | "and" | "or" | "xor" | "shiftl" | "lshiftr" |
			"ashiftr" | "add" | "sub" | "mult" | "divfloor" | "exp" | "mod") +
			number
	expression = call | ref | deref | literal | operation
	segment = "segment" + name
	function = "function" + OneOrMore(Forward(statement)) + "endfunction"
	pop = "pop" + Optional(number) + Optional(address)
	statement << (segment | function | pop | ifstmt)

	return statement.parseString(code)

def list_instructions(code):
	p = parse(code)
	print(p)

