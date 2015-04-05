
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

class terminal(int):
    pass

# Grammar is determined by G = (N, T, P, S)
# The constructor builds this set from the input, and transforms N,P into a
# canonical form.  The hash value of grammar objects are hence based on N, T,
# and P.
class grammar:
    def utf8map(value):
        if type(value) == str: return terminal(ord(value))
        elif type(value) == terminal: return chr(value)
        else: raise Exception('Invalid utf8 character')
    def __init__(self, input, start='S', terminal_map=utf8map):
        self.N = set()
        self.T = set()
        self.P = []
        self.S = start
        tokens = ('EMPTYSTRING', 'EQUALS', 'NEWLINE', 'NONTERMINAL', 'TERMINAL')
        t_ignore_space = r'\ +'
        t_EQUALS = r'='
        t_EMPTYSTRING = r'ε'
        def t_NEWLINE(t):
            r'\\?\n'
            t.lexer.lineno += 1
            if t.value == '\\\n': return None
            return t
        ### TODO: replace terminals and non-terminals with symbol container so
        ### that a symbol can be replaced en-masse with another (and a terminal
        ### can be replaced with a non-terminal and vice-versa)
        def t_TERMINAL(t):
            r'(".")|(0x[0-9a-f]+)'
            if t.value[0] == '"': t.value = terminal_map(t.value[1])
            else: t.value = terminal(t.value[2:], 16)
            self.T.add(t.value)
            return t
        def t_NONTERMINAL(t):
            r'[^\W\d][\w-]*'
            self.N.add(t.value)
            return t
        def p_S(p):
            '''S : production S
                 | NEWLINE S
                 | production
                 | NEWLINE'''
        def p_production(p):
            '''production : string EQUALS string NEWLINE
                          | string EQUALS EMPTYSTRING NEWLINE'''
            p[0] = (p[1], p[3])
            self.P.append(p[0])
        def p_string(p):
            '''string : NONTERMINAL string
                      | TERMINAL string
                      | NONTERMINAL
                      | TERMINAL'''
            if len(p) == 3: p[0] = [ p[1] ] + p[2]
            else: p[0] = [ p[1] ]
        lexer = lex.lex()
        parser = yacc.yacc()
        parser.parse(input)
        self.canonical()
    def __repr__(self):
        out = ''
        for p in self.P:
            lhs = ' '.join([ s if type(s) == str else str(s) for s in p[0] ])
            rhs = ' '.join([ s if type(s) == str else str(s) for s in p[1] ])
            out += lhs + ' = ' + rhs + '\n'
        return ''
    # The wiki page on 'Kuroda normal form' mentions that any unrestricted language
    # is weakly equivalent to one where all rules are of the form:
    #   1.  AB = CD
    #   2.  A = BC
    #   3.  A = a
    #   4.  A = ε
    # If rule type 4. is ommitted, we have a normal form for context-sensitive
    # languages, whereas if only rule 1. is ommitted, we gain a normal form for
    # context-free languages.
    def canonical(self):
        # Replace all terminals with non-terminals
        T = sorted(self.T)
        T = dict(zip(T, [ "'" + str(t) for t in T ]))
        for lhs, rhs in self.P:
            for i in range(len(lhs)):
                if lhs[i] in T: lhs[i] = T[lhs[i]]
                if rhs[i] in T: rhs[i] = T[rhs[i]]
            print(lhs + rhs)
        # Find all productions in P of the form 1, 2, 3, or 4

