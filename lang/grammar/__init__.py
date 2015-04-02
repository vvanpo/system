
# Chomsky hierarchy
#   Unrestricted
#   Context-sensitive
#   Context-free
#   Regular
# Production notation:
# This notation uses itself to describe itself.  Informally, it can describe any
# grammar using regular productions (a = b) along with postfix operators to shorthand
# options ("|"), zero-or-one ("?"), zero-or-more ("*"), one-or-more ("+"), specific
# range of repeats ("{M,N}", where M and N are integers, or one of them is missing to
# denote 0 or infinity, respectively.  M without a comma denotes an exact number of
# repeats), and grouping (using parentheses).  Symbol concatenation is achieved with
# whitespace, and terminals are indicated with double quotes (although a character set
# must be defined beforehand).
#   S = (space | newline)* production+ (space | newline)*
#     production = string | empty-string "=" string | empty-string newline+
#       string = space? group | concat space?
#         group = ("(" space? group | concat space? ")") | symbol
#                   | (group (zero-or-one | zero-or-more | one-or-more | range))
#           symbol = terminal | non-terminal
#           zero-or-one = "?"
#           zero-or-more = "*"
#           one-or-more = "+"
#           range = "{" range-exact | (range-low "," range-high?) | ("," range-high) "}"
#             range-exact = digit+
#             range-low = digit+
#             range-high = digit+
#         concat = group | option (space group | option)+
#           option = group (space? "|" space? group)+
#       space = " "+
#         ###escaped-newline = "\" newline
#       non-terminal = (letter | digit) (letter | digit | "-")*
#       terminal = ('"' character '"') | hex
#         hex = ("0" .. "9" | "a" .. "f")+
#       empty-string = space? "ε" space?

from ply import lex, yacc

class grammar:
    def __init__(self, input, character_set='utf-8'):
        tokens = ('EMPTYSTRING', 'EQUALS', 'NEWLINE', 'NONTERMINAL', 'TERMINAL')
        t_ignore_space = r'\ +'
        t_EQUALS = r'='
        t_NONTERMINAL = r'\w[\w-]*'
        t_EMPTYSTRING = r'ε'
        def t_NEWLINE(t):
            r'\\?\n'
            t.lexer.lineno += 1
            if t.value == '\\\n': return None
            return t
        def t_TERMINAL(t):
            r'(".")|([0-9a-f]+)(?!\w)'
            if t.value[0] == '"': t.value = ord(t.value[1])
            else: t.value = int(t.value, 16)
            return t
        lexer = lex.lex()
        lexer.input(input)
        for t in lexer: print(t)

