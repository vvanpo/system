import asm, re

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
#       integer =   { "0" .. "9" | "a" .. "f" }-
class data:
    # Captured groups:
    r = re.compile(r'\s*?(?:'
            + r'{(.+?)}|'       # 1 label
            + r'([0-9a-f]+)|'   # 2 integer
            + r'('              # 3 string
              + r'r"(?:\\"|.)*?(?<![^\\]\\)"|' # raw string
              + r'(?:"(?:(?<!\\)(?:\\\\)*"|.)*?(?<!\\)(?:\\\\)*")' # regular string
            + r')'              # /string
            + r')(?:\s*,)?|(?:\s*$)')
    @classmethod
    def from_string(cls, stmt):
        value = bytearray()
        # Dictionary of index -> label-expression
        labels = {}
        def parse_string(s):
            nonlocal value, labels
            # TODO: make encoding an architecture option
            add_val = lambda x: value.extend(x.encode('utf-8'))
            raw = False
            if s[0] == 'r':
                s = s[1:]
                raw = True
            s = s[1:-1] # string double quotes
            if not raw:
                s = re.split(r'((?:\\\\)+)|\\{(.+?)}|\\([abfnrtv"])|\\x([0-9a-f]{2})|\\u([0-9a-f]{4})|\\U([0-9a-f]{8})', s)
                for i in range(0, len(s)-1, 7):
                    m = s[i:i+7]
                    add_val(m[0])
                    if m[1]: add_val(m[1][::2])
                    if m[2]:
                        if len(value) not in labels: labels[len(value)] = []
                        labels[len(value)].append('"' + m[2] + '"')
                    if m[3]: add_val('\a\b\f\n\r\t\v"'['abfnrtv"'.index(m[3])])
                    if m[4]: value.extend(bytes.fromhex(m[4]))
                    if m[5]: add_val(chr(int(m[5], 16)))
                    if m[6]: add_val(chr(int(m[6], 16)))
            else:
                s = re.split(r'\\(["])', s)
                for i in range(0, len(s)-1, 2):
                    add_val(s[i])
                    if s[i+1]: add_val('"')
            add_val(s[-1])
        s = cls.r.split(stmt)
        for i in range(0, len(s)-1, 4):   # number_matches = (s - 1)/4
            m = s[i:i+4]
            if m[0]: raise Exception("Invalid data statement: " + stmt)
            if m[1]:
                if len(value) not in labels: labels[len(value)] = []
                labels[len(value)].append(m[1])
            if m[2]:
                num = int(m[2], 16)
                align = 1   # TODO: add option in data sections
                l = (num.bit_length() // (align*8)) + 1
                value.extend(num.to_bytes(l, 'little'))
            if m[3]: parse_string(m[3])
        return cls(value, labels)
    def __init__(self, value, labels):
        self.value = value
        self.labels = labels
    def __str__(self):
        return str(self.value)
    __repr__ = __str__

# Register the 'data' architecture name with ..asm package
asm.architecture.register('data', data)

