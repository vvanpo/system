
from pyparsing import *

def parse(code):
	ParserElement.setDefaultWhitespaceChars(" \t\r")
	decimal = Word(nums)
	decimal.setParseAction(lambda t: int(t[0]))
	hex = ("0x" + Word(hexnums)).leaveWhitespace()
	hex.setParseAction(lambda t: int(t[1], 16))
	number = decimal | hex
	address = Optional(Optional("segment") + ("ref" | "val") + oneOf("ip sp fp"))
	literal = "literal" + number
	load = "load" + Optional("bytes") + address
	store = "store" + Optional("bytes") + address
	call = "call" + address
	ifzero = "if zero" + address
	op = oneOf("not and or xor shiftl lshiftr ashiftr add sub mult floordiv exp mod") + Optional("bytes")
	simple = oneOf("open close return")

	instr = lineStart.suppress() + (literal | load | store | call | jump | ifzero | op | simple) + lineEnd.suppress()
	instr.setWhitespaceChars(" \t\r\n")
	asmlang = stringStart + OneOrMore(instr) + stringEnd

	instrs = []
	def action_instr(s, loc, toks):
		instr.append(toks.asList())
	instr.setParseAction(action_instr)
	asmlang.parseString(code)
	return instr

# sets up a call frame
def frame(addr):
	f = """
		load val fp
		load val sp
		store val fp
		literal """ + str(addr) + """
		call
		store val fp
		"""
	return f
