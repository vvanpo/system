import asm

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
    @classmethod
    def from_string(cls, stmt):
        def parse_string(stmt, value, labels):
            raw = False
            ascii = False
            i = stmt.find('"')
            if i > 0:
                if stmt[0] == 'r': raw = True
                if stmt[0] == 'b' or stmt[1] == 'b': ascii = True
            while True:
                if i == -1: raise Exception("Unending string")
                i = stmt.find('"', i+1)
                if stmt[i-1] != "\\": break
            while True:
                j = stmt.find(r'\{')
                if j == -1: break
                k = stmt.find('}', j)
                labels[j] = stmt[j+2:k]
                stmt = stmt[:j] + stmt[k+1:]
            value.extend(stmt[:i].encode())
            return stmt[i+1:].strip()
        def parse_integer(stmt, value):
            i = stmt.find(',')
            if i == -1: i = len(stmt)
            value.extend(stmt[:i].encode())
            return stmt[i+1:].strip()
        value = bytearray()
        # Dictionary of index -> label-expression
        labels = {}
        stmt = stmt.strip()
        while stmt:
            # Check for string
            if stmt[0] == '"' or len(stmt) > 2 and (stmt[1] == '"' or stmt[:2] == 'rb"'):
                stmt = parse_string(stmt, value, labels)
                continue
            stmt = parse_integer(stmt, value)
        return cls(bytes(value), labels)
    def __init__(self, value, labels):
        self.value = value
        self.labels = labels
    def __str__(self):
        return str(self.value)
    __repr__ = __str__

# Register the 'data' architecture name with ..asm package
asm.architecture.register('data', data)

