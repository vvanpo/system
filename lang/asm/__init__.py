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
        # List of labels
        self.labels = []
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
                section = self.format.new_section(w[1], architecture.new(a[0]), option_string)
                continue
            if not self.format.sections: raise Exception("No section declarations")
            ## Label
            # In the case of data statements, it's possible that a ':' could be
            # found in the data string instead, so we need to make sure it isn't
            # a string by checking for '"'
            if w[0][-1] == ':' and w[0][0] != '"':
                name = w[0][:-1]
                l = l[l.index(':')+1:]
                del w[0]
                lbl = label(name, section)
                if lbl in self.labels:
                    raise Exception("Duplicate label: " + str(section) + "." + name)
                self.labels.append(lbl)
                if not w: continue
            ## Instruction
            section.add_statement(l)
        self.format.calculate_addr()

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

class label:
    @staticmethod
    def expression(e):
        pass
    def __init__(self, name, section):
        if not re.match(r"[\w-]+$", name):
            raise Exception("Invalid label: " + str(section) + "." + name)
        self.name = name
        self.section = section
    def __eq__(self, other):
        if self.name == other.name and self.section == other.section: return True
        return False
    def __str__(self):
        return self.name

# Import all modules in . directory so they will register with architecture
for m in [ '.' + p for p in os.listdir(*__path__) if p[0] != '.' and p[0] != '_' ]:
    importlib.import_module(m, __name__)

