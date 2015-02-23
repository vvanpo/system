import re, importlib, os

# Format of all assembly files is simple:
#   assembly =          [ newline ], architecture, { section }-
#     ; architecture (x86-64, i686, etc.) [ (elf, bin, etc.) ]
#     architecture =    "architecture", space, identifier, [ space, identifier ], newline
#     section =         "section", space, section-name, newline, instruction-list
#       section-name =  identifier
#       instruction-list = 
#                       { ( label, [ newline ], instruction ) | instruction }-
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
        # Set of sections
        self.sections = set()
        # Set of labels
        self.labels = set()
        line_no = 0
        for l in source.splitlines():
            line_no += 1
            w = l.split()
            if len(w) == 0: continue
            ## Architecture
            if not hasattr(self, "architecture"):
                if not w[0] == "architecture" or len(w) < 2:
                    raise Exception("Missing architecture declaration")
                self.architecture = architecture.get_instance(w[1])
                # TODO: enumerate different formats
                if len(w) > 2: self.format = w[2:]
                continue
            ## Section
            if w[0] == "section":
                s = section(w[1])
                if s in self.sections: raise Exception("Duplicate section: " + s)
                self.sections.add(s)
                continue
            if not self.sections: raise Exception("No section declarations")
            ## Label
            if w[0][-1] == ":":
                n = label(s, w[0][:-1], len(s.instructions) + 1)
                if n in self.labels: raise Exception("Duplicate label: " + n)
                self.labels.add(n)
                del w[0]
                if not w: continue
            ## Instruction
            # Instructions aren't assembled until after label addresses can be
            # calculated.  This can require multiple passes.
            ###instruction = self.architecture.instruction.from_string(" ".join(w))
            s.instructions.append(" ".join(w))

class architecture:
    names = {}
    @classmethod
    def register(cls, name, instruction_class):
        self = cls(instruction_class)
        self.name = name
        cls.names[name] = self
    def __init__(self, instruction_class):
        self.instruction = instruction_class
    @classmethod
    def get_instance(cls, name):
        if name not in cls.names:
            raise Exception("Invalid architecture name: " + name)
        return cls.names[name]

class section(str):
    def __new__(cls, name):
        self = super().__new__(cls, name)
        if not re.match(r"[\w-]+$", self): raise Exception("Invalid section name")
        self.instructions = []
        return self

class label(str):
    def __new__(cls, section, name, index):
        self = super().__new__(cls, section + "." + name)
        if not re.match(r"[\w-]+\.[\w-]+$", self): raise Exception("Invalid label")
        self.section = section
        self.name = name    # name + section describe a unique label
        self.index = index  # index into section matching instruction pointed to
        return self

# Import all modules in . directory so they will register with architecture
modules = [ '.' + p for p in os.listdir(*__path__) if p[0] != '.' and p[0] != '_' ]
for m in modules:
    importlib.import_module(m, __name__)

