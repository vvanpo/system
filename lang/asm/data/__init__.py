import asm, re

class label_expression:
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
    # label_expressions are stored as an AST unless there are no labels in the expression,
    # in which case it calculates the value
    # TODO: clean this up, and put in some error-checking
    def __init__(self, string, encoding):
        self.encoding = encoding
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
                    + r'((?<!\w)(?:[0-9a-f]+h|[0-9]+)(?!\w))|' # integer
                    + r'((?<!\w)\w[\w.-]*)|'   # label
                    + r'(\*\*|\*|/|%|\+|-)|'    # operator
                    + r'(\()|'                  # open-bracket
                    + r'(\))'                   # close-bracket
                    + r')\s*'
                    , string)
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
        self.expr = parse(tokens)
        ###TEST########################################
        print(self, self.resolve({
            "_end": 0, "testlabel": 1, "data.end": 2, "code.end": 3, "data.hello": 4,
            "hello": 5, "code.testlabel": 6, "a": 7,
            }))
        ##################################################
    def __repr__(self):  return repr(self.expr)
    def resolve(self, labels):
        # labels should be a dict of names -> addresses
        def match_label(string):
            nonlocal labels
            if string not in labels:
                raise Exception("Non-existing label in label expression")
            return labels[string]
        def recurse(tup):
            if type(tup) == str: return match_label(tup)
            if type(tup) == int: return tup
            if type(tup) != tuple or len(tup) != 3: raise Exception("Invalid label expression")
            return tup[1](recurse(tup[0]), recurse(tup[2]))
        return recurse(self.expr)

# TODO: parts of this should be factored out into the asm module
class statement(bytearray):
    def __init__(self, encoding, *args):
        super().__init__(*args)
        # Dictionary of index -> label-expression
        self.label_exprs = {}
        self.encoding = encoding
    def __repr__(self):
        return super().__repr__()[10:-1] + ' -> ' + str(self.label_exprs)
    def add_str(self, string): self.extend(string.encode(self.encoding))
    def add_label_expr(self, string, encoding=False):
        encoding = self.encoding if encoding else None
        label_expr = label_expression(string, encoding)
        if len(self) not in self.label_exprs:
            self.label_exprs[len(self)] = []
        self.label_exprs[len(self)].append(label_expr)
    def length(self):
        # length() returns -1 if the statement is as-yet unresolved
        if self.label_exprs: return -1
        return super().__len__()
    def resolve_labels(self, labels):
        # labels should be a dict of names -> addresses
        for i, l_list in self.labels.copy().items():
            for j in range(len(l_list)):
                del l_list[j]
            if not l_list: del self.labels[i]

# data is a simple 'architecture' for storing strings and numbers and the like
# in binaries
# Syntax for various kinds of data:
#   statement =     { value, ",", space }, value
#     value =       string | integer
#       ; string syntax is very similar to python3 strings
#       string =    [ prefix ], '"', { character }, '"'
#         prefix =  "r"
#         character =   ? any printable character except for newline ?
#         ; escapes are included for the usual C escapes, \xxx bytes, \uxxxx, \Uxxxxxxxx,
#         ; and label arithmetic embedded in strings using \{}
#       ; label arithmetic is also available to specify integers, and doesn't always
#       ; require actual labels, e.g. {2**8 - 1}
#       integer =   { "0" .. "9" }- | ( { "0" .. "9" | "a" .. "f" }-, [ "h" ] )
class data(asm.architecture):
    _stmt_regex = re.compile(r'\s*?(?:'
            + r'{(.+?)}|'                   # 1 label
            + r'([0-9a-f]+h|[0-9]+)|'       # 2 integer
            + r'('                          # 3 string
              + r'r"(?:\\"|.)*?(?<![^\\]\\)"|' # raw string
              + r'(?:"(?:(?<!\\)(?:\\\\)*"|.)*?(?<!\\)(?:\\\\)*")' # regular string
            + r')'
            + r')(?:\s*,)?|(?:\s*$)')
    _str_regex = re.compile(r'((?:\\\\)+)|' # 1 even number of \
            + r'\\{(.+?)}|'                 # 2 label expression
            + r'\\([abfnrtv"])|'            # 3 C escape character
            + r'\\x([0-9a-f]{2})|'          # 4 byte number
            + r'\\u([0-9a-f]{4})|'          # 5 unicode codepoint (16-bit)
            + r'\\U([0-9a-f]{8})')          # 6 unicode codepoint (32-bit)
    def __init__(self, option_string):
        self.align = 1
        self.encoding = 'utf-8'
        self.endian = 'little'
        m = re.match(r'(?:\s*(?:align=([0-9]+)|'
                    + r'encoding=(utf-8)|endian=(little|big)))+', option_string)
        if m and m.group(1): self.align = int(m.group(1))
        if m and m.group(2): self.encoding = m.group(2)
        if m and m.group(3): self.endian = m.group(3)
        self.statements = []
    def __repr__(self):
        return 'data:\n' + ''.join([ '   ' + repr(s) + '\n' for s in self.statements ])
    def add_statement(self, string):
        value = statement(self.encoding)
        def parse_string(s):
            nonlocal value
            raw = False
            if s[0] == 'r':
                s = s[1:]
                raw = True
            s = s[1:-1] # string double quotes
            if not raw:
                s = self.__class__._str_regex.split(s)
                for i in range(0, len(s)-1, 7):
                    m = s[i:i+7]
                    value.add_str(m[0])
                    if m[1]: value.add_str(m[1][::2])
                    if m[2]: value.add_label_expr(m[2], True)
                    if m[3]: value.add_str('\a\b\f\n\r\t\v"'['abfnrtv"'.index(m[3])])
                    if m[4]: value.extend(bytes.fromhex(m[4]))
                    if m[5]: value.add_str(chr(int(m[5], 16)))
                    if m[6]: value.add_str(chr(int(m[6], 16)))
            else:
                s = re.split(r'\\(["])', s)
                for i in range(0, len(s)-1, 2):
                    value.add_str(s[i])
                    if s[i+1]: value.add_str('"')
            value.add_str(s[-1])
        def parse_int(s):
            nonlocal value
            if s[-1] == 'h': num = int(s[:-1], 16)
            else: num = int(s)
            l = self.align * (num.bit_length() // (self.align * 8)) + self.align
            value.extend(num.to_bytes(l, self.endian))
        s = self.__class__._stmt_regex.split(string)
        for i in range(0, len(s)-1, 4):   # number_matches = (s - 1)/4
            m = s[i:i+4]
            if m[0]: raise Exception("Invalid data statement: " + string)
            if m[1]: value.add_label_expr(m[1])
            if m[2]: parse_int(m[2])
            if m[3]: parse_string(m[3])
        if s[-1]: raise Exception("Invalid data value: " + s[-1])
        self.statements.append(value)
        #value.resolve_labels({})

# Register the 'data' architecture name with ..asm package
asm.architecture.register('data', data)

