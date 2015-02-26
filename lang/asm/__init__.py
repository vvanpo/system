import re, importlib, os
from asm.formats import bin_format

# Format of all assembly files is simple:
#   assembly =          [ newline ], format, [ default_architecture ], { section }-
#     ; format (elf, bin, etc.)
#     format = "format", space, identifier, newline
#     default_architecture = "architecture", space, architecture, newline
#       ; architecture (x86-64, i686, etc.)
#       architecture =    identifier
#     section =         "section", space, section-name, [ architecture ], [ format-options ],
#                       newline, instruction-list
#       section-name =  identifier
#       format-options = ? string ?     ; dependent on format
#       instruction-list = 
#                       { ( label, newline ) | ( label, space, instruction ) | instruction }-
#         label =       identifier, ":"
#       ; TODO: need to fix this to accomodate any architecture's assembly
#       ; Might want to be able to parse labels in operands, and do label arithmetic
#       instruction =   ? some string ?, newline
#     identifier =  letter, { letter | digit }
#       letter =    "_" | "-" | ? a Unicode code point classified as "Letter" ?
#       digit =     ? a Unicode code point classified as "Decimal Digit" ?
#     space =       { ? the Unicode code points U+0020 and U+0009 ? }-
#     newline =     { ? the Unicode code point U+000a ? }-
class assembly:
    def __init__(self, source):
        self.default_architecture = ('ir', "")
        lines = source.splitlines()
        line_no = 0
        for l in lines:
            line_no += 1
            w = l.split()
            if not w: continue
            ## Format
            if not hasattr(self, 'format'):
                if not w[0] == 'format' or len(w) < 2:
                    raise Excepion("Missing format declaration")
                self.format = bin_format.new(w[1])
                continue
            ## Architecture
            if w[0] == "architecture":
                self.default_architecture = (w[1], " ".join(w[2:]))
                continue
            ## Section
            if w[0] == "section":
                a = self.default_architecture
                option_string = ""
                if len(w) > 2:
                    try:
                        # TODO: arch options
                        a = (w[2],"")
                        format_options = " ".join(w[3:])
                    except Exception: format_options = " ".join(w[2:])
                section = self.format.add_section(w[1], architecture.new(a[0]), option_string)
                continue
            if not self.format.sections: raise Exception("No section declarations")
            ## Label
            # In the case of data statements, it's possible that a ':' could be
            # found in the data string instead, so we need to make sure it isn't
            # a string by checking for '"'
            # BUG: TODO: doesn't check for raw string
            if w[0][-1] == ':' and w[0][0] != '"':
                section.add_label(w[0][:-1])
                l = l[l.index(':')+1:]
                if len(w) == 1: continue
            ## Instruction
            section.add_statement(l)
        self.format.assemble()
        print(self.format)

class architecture:
    names = {}
    @classmethod
    def register(cls, name, architecture):
        cls.names[name] = architecture
    @classmethod
    def new(cls, name, options=""):
        if name not in cls.names:
            raise Exception("Invalid architecture name: " + name)
        return cls.names[name](options)

# Import all modules in . directory so they will register with architecture
for m in [ '.' + p for p in os.listdir(*__path__) if p[0] != '.' and p[0] != '_' ]:
    importlib.import_module(m, __name__)

