from . import ir, x86
import re

# Format of all assembly files is simple:
#   assembly = 		    [ newline ], architecture, { section }-
#     ; architecture (x86-64, i686, etc.) [ (elf, bin, etc.) ]
#     architecture =	"architecture", space, identifier, [ space, identifier ], newline
#     section =		    "section", space, section-name, newline, instruction-list
#       section-name =	identifier
#       instruction-list = 
#                       { ( label, [ newline ], instruction ) | instruction }-
#         label = 	    identifier, ":"
#       ; TODO: need to fix this to accomodate any architecture's assembly
#       ; Might want to be able to parse labels in operands, and do label arithmetic
#       instruction =	? some string ?, newline
#     identifier =	letter, { letter | digit }
#       letter =	"_" | "-" | ? a Unicode code point classified as "Letter" ?
#       digit =		? a Unicode code point classified as "Decimal Digit" ?
#     space =       { ? the Unicode code points U+0020 and U+0009 ? }-
#     newline =		{ ? the Unicode code point U+000a ? }-
class assembly:
    def __init__(self, source):
        def architecture_add(name, format=None):
            for a in (ir, x86):
                if name in a.arch_strings:
                    self.architecture = a
                    if not format: format = a.default_format
                    break
            if not hasattr(self, "architecture"):
                raise Exception("Invalid architecture declaration")
            # TODO: enumerate different formats
            self.format = format
        # Dictionary of section names -> [instruction]
        self.sections = {}
        # Dictionary of label names -> instructions
        self.labels = []
        line_no = 0
        for l in source.splitlines():
            line_no += 1
            w = l.split()
            if len(w) == 0: continue
            ## Architecture
            if not hasattr(self, "architecture"):
                if not w[0] == "architecture" or len(w) < 2:
                    raise Exception("Missing architecture declaration")
                architecture_add(*w[1:])
                continue
            ## Section
            if w[0] == "section":
                s = section(w[1])
                if s in self.sections: raise Exception("Double section name " + s)
                self.sections[s] = []
                current_section = s
            if not self.sections: raise Exception("No section declarations")
            ## Label
            if w[0][-1] == ":":
                name = label(current_section, w[0][:-1])
                if name in self.labels: raise Exception("Duplicate label")
                self.sections[current_section].append(name)
                del w[0]
                if not w: continue
            ## Instruction
            instruction = self.architecture.instruction.from_string(" ".join(w))
            self.sections[current_section].append(instruction)

class section(str):
    def __new__(cls, *args):
        if not re.match(r"[\w-]*$", super().__new__(cls, *args)):
            raise Exception("Invalid section name")

class label(str):
    def __new__(cls, section, name):
        if not re.match(r"[\w-]*$", self): raise Exception("Invalid label")
        ###
