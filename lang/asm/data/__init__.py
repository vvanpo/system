import asm, re

# TODO: parts of this should be factored out into the asm module
class statement(bytearray):
    def __init__(self, encoding, *args):
        super().__init__(*args)
        # Dictionary of index -> label-expression
        self.labels = {}
        self.encoding = encoding
    def __repr__(self):
        return super().__repr__()[10:-1] + ' -> ' + str(self.labels)
    def add_str(self, string): self.extend(string.encode(self.encoding))
    def add_label(self, label):
        if len(self) not in self.labels:
            self.labels[len(self)] = []
        self.labels[len(self)].append(label)
    def length(self):
        # length() returns -1 if the statement is as-yet unresolved
        if self.labels: return -1
        return super().__len__()
    def resolve_labels(self, labels):
        # labels should be a dict of names -> addresses
        for i, l_list in self.labels.copy().items():
            for j in range(len(l_list)):
                encode = False
                l = l_list[j]
                if l[0] == '"':
                    encode = True
                    l = l[1:-1]
                #del l_list[j]
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
                    if m[2]: value.add_label('"' + m[2] + '"')
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
            if m[1]: value.add_label(m[1])
            if m[2]: parse_int(m[2])
            if m[3]: parse_string(m[3])
        if s[-1]: raise Exception("Invalid data value: " + s[-1])
        self.statements.append(value)
        #value.resolve_labels({})

# Register the 'data' architecture name with ..asm package
asm.architecture.register('data', data)

