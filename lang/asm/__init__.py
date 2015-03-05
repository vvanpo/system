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
#     ; TODO: identifiers need to be distinguishable from integers
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
            if w[0][-1] == ':':
                section.add_label(w[0][:-1])
                l = l[l.index(':')+1:]
                if len(w) == 1: continue
            ## Instruction
            section.add_statement(l)
        self.format.assemble()
    def __repr__(self): return repr(self.format)
    @staticmethod
    # Label expressions are embedded within statements using {} as delimiters
    # expression =          add_sub
    #   add_sub_expr =      ( mult_div, "+", mult_div )
    #                       | ( [ mult_div ], "-", mult_div ) | mult_div
    #   mult_div =          ( exp, "*", exp ) | ( exp, "/", exp ) | exp
    #   exp =               ( bracket, "**", bracket ) | bracket
    #   bracket =           ( "(", expression, ")" ) | integer | label
    #   integer =           decimal | hex
    #     decimal =         { "0" .. "9" }-
    #     hex =             { "0" .. "9" | "a" .. "f" }-, "h"
    #   label =             ? identifier as above ?
    # label_expression() returns an AST unless there are no labels in the expression,
    # in which case it calculates the value
    # TODO: clean this up, and put in some error-checking
    def label_expression(e):
        add = int.__add__
        sub = int.__sub__
        mul = int.__mul__
        div = int.__floordiv__
        mod = int.__mod__
        pow = int.__pow__
        ops = ((add, sub), (mul, div, mod), (pow,))
        def parse(tokens):
            nonlocal ops
            # TODO: this return statement allows unary '-' operator to exist, but we need to disallow unary versions of the other operators
            if not tokens: return 0
            if len(tokens) == 1: return tokens[0]
            for i in range(len(tokens)):
                if tokens[i] == '(':
                    queue = 0
                    j = 1
                    for t in tokens[i+1:]:
                        if t == '(': queue += 1
                        if t == ')':
                            if queue == 0:
                                return parse(tokens[:i] + [ parse(tokens[i+1:i+j]) ] + tokens[i+j+1:])
                            queue -= 1
                        j += 1
            for group in ops:
                for i in range(0, -len(tokens), -1):
                    if tokens[i] in group:
                        left = parse(tokens[:i])
                        right = parse(tokens[i+1:])
                        if type(left) == int and type(right) == int:
                            return tokens[i](left, right)
                        return (left, tokens[i], right)
        s = re.split(r'\s*(?:'
                    + r'((?<!\w)(?:[0-9a-f]+h|[0-9]+)(?!\w))|'   # integer
                    + r'((?<!\w)[\w][\w-]*)|'   # label
                    + r'(\*\*|\*|/|%|\+|-)|'    # operator
                    + r'(\()|'                  # open-bracket
                    + r'(\))'                   # close-bracket
                    + r')\s*'
                    , e)
        integer = lambda s: int(s[:-1], 16) if re.match(r'[0-9a-f]+h', s) else int(s) if re.match(r'[0-9]+', s) else None
        other = lambda s: s
        operator = lambda s: add if s == '+' else sub if s == '-' else mul if s == '*' \
                        else div if s == '/' else mod if s == '%' else pow if s == '**' \
                        else None
        tokens = []
        for i in range(0, len(s)-1, 6):
            if s[i] or s[-1]: raise Exception("Invalid label expression: " + e)
            tags = (integer, other, operator, other, other)
            for j in range(5):
                if s[i+j+1]: tokens.append(tags[j](s[i+j+1]))
        return parse(tokens)

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

