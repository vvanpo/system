import asm, re

# data is a simple 'architecture' for storing strings and numbers and the like
# in binaries
# Syntax for various kinds of data:
#   statement =     { value, ",", space }, value
#     value =       string | integer
#       ; string syntax is very similar to python3 strings
#       string =    [ prefix ], '"', { character | escape }, '"'
#         prefix =  "r" | "b" | "rb"        ; raw, ascii, raw & ascii
#         character =   ? any printable character except for \, newline, or double-quote ?
#         ; escapes are included for the usual C escapes, \xxx bytes, \uxx characters,
#         ; and label arithmetic embedded in strings using \{}
#         escape =  "\", ? any printable character ?
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
            raw = False
            if s[0] == 'r':
                s = s[1:]
                raw = True
            s = s[1:-1]
            if not raw:
                m = re.split(r'(\\{.+?})', s)
                for j in range(0, len(m), 2):
                    # TODO: make encoding a per-line option
                    value.extend(m[j].encode('utf-8'))
                    if j < len(m)-1:
                        if len(value) not in labels: labels[len(value)] = []
                        labels[len(value)].append(m[j+1])
        s = cls.r.split(stmt)
        num_groups = 4
        for i in range(0, len(s)-1, num_groups):   # number_matches = (s - 1)/num_groups
            m = s[i:i+num_groups]
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
        return cls(bytes(stmt.encode()), labels)
    def __init__(self, value, labels):
        self.value = value
        self.labels = labels
    def __str__(self):
        return str(self.value)
    __repr__ = __str__

# Register the 'data' architecture name with ..asm package
asm.architecture.register('data', data)

