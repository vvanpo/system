
# Format of all assembly files is simple:
#   assembly = 		format, section+
#     format =		"format", format-name, architecture, newline
#       format-name =	identifier
#	architecture =	identifier
#     section =		"section", section-name, newline, 
#       section-name =	identifier
#       instruction =	[ label ], 
#         label = 	identifier, ":"
#     identifier =	letter, { letter | digit }
#       letter =	"_" | ? a Unicode code point classified as "Letter" ?
#       digit =		? a Unicode code point classified as "Decimal Digit" ?
#     newline =		? the Unicode code point U+000a ?
class assembly:
	def __init__(self, source):
		line_no = 0
		for l in source.splitlines():
			line_no++
			if l.strip() == "": continue

class label:
	pass

class instruction:
	def __init__(self, string):
		pass
