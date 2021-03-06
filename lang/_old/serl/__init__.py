from ply import lex, yacc

# An ordered map has all the same operations as a regular list, but items that are
# pairs can be indexed by their key (which returns their value), or their numerical
# index (which returns the key-value pair)
class omap(list):
    def __init__(self, *args):
        super().__init__(*args)

class pair:
    def __init__(self, key, value):
        self.key = key
        self.value = value
    def __repr__(self):
        return "p{ " + repr(key) + " -> " + repr(value) + " }"

# serl: Serialization Language
# A data format that can map to sets, dicts, lists, etc. and is human-readable.
# The character set, for now, UTF-8 by default.
# For now, indent tokens are handled outside the grammar, until I can figure out
# how to write a proper context-sensitive grammar, and corresponding parser.
# Reference for CSGs: http://www.diku.dk/hjemmesider/ansatte/henglein/papers/chomsky1959.pdf
#   S = omap? ENDOFFILE
#     ENDOFFILE = "——" NEWLINE
#     omap = plain-item+ | dashed-item+
#       plain-item = pair | scalar
#       dashed-item = "- " plain-item | (complex-key (":" value)?) | embedded-omap
#     embedded-omap = (dashed-item (INDENT dashed-item+ DEDENT)?)
#                       | (plain-item INDENT plain-item+ DEDENT)
#     pair = key ":" value
#       key = scalar-inline
#       value = scalar
#     scalar-inline = 
#     complex-key = 
class file:
    def __init__(self):
        self.root = omap()
    def __repr__(self):
        return str(self.root)
    @classmethod
    def from_string(cls, string):
        self = cls()
        tokens = self.tokenize(string)
        return self
    @classmethod
    def from_variable(cls, var):
        pass
    def tokenize(self, string):
        # current indent position
        indent = 0
        tokens = (
            'ENDOFFILE',
            'NEWLINE',
            'INDENT',
            'DEDENT',
            )
        t_ignore_z = r'.+?'
        def t_INDENT(t):
            r'\n(\ \ )*(?=[^ \n])'
            nonlocal indent
            current = len(t.value)//2
            if current != indent:
                diff = current - indent
                if diff < 0: t.type = 'DEDENT'
                if abs(diff) > 1:
                    indent += diff//abs(diff)
                    t.value = indent
                    t.lexer.lexpos = t.lexpos   # Re-lex token
                    return t
                t.value = indent = current
            else: t.type = 'NEWLINE'
            t.lexer.lineno += 1
            return t
        def t_NEWLINE(t):
            r'\n(\ \ )*'
            t.lexer.lineno += 1
            return t
        def t_ENDOFFILE(t):
            r'——\n'
            nonlocal string
            t.lexer.lexpos = len(string)    # stop lexing
            return t
        lexer = lex.lex()
        lexer.input(string)
        for t in lexer:
            print(t)

def parse(source):
    return file.from_string(source)

def dump(var):
    return file.from_variable(var)

