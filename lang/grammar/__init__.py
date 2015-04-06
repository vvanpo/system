
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

# Grammar is determined by G = (N, T, P, S)
# The constructor builds this set from the input, and transforms N,P into a
# canonical form.  The hash value of grammar objects are hence based on N, T,
# and P.
class grammar:
    def utf8map(value):
        if type(value) == str: return ord(value)
        elif type(value) == int: return chr(value)
        else: raise Exception('Invalid utf8 character')
    def __init__(self, input, S='S', terminal_map=utf8map):
        self.P = []
        self.S = S
        # symbols is the union of N and T
        self.symbols = [S]
        tokens = ('EMPTYSTRING', 'EQUALS', 'NEWLINE', 'NONTERMINAL', 'TERMINAL')
        t_ignore_space = r'\ +'
        t_EQUALS = r'='
        def t_EMPTYSTRING(t):
            r'ε'
            return t
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
            else: t.value = int(t.value[2:], 16)
            if t.value not in self.symbols:
                self.symbols.append(t.value)
            t.value = self.symbols.index(t.value)
            return t
        def t_NONTERMINAL(t):
            r'[^\W\d][\w-]*'
            if t.value not in self.symbols:
                self.symbols.append(t.value)
            t.value = self.symbols.index(t.value)
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
            if p[3] == 'ε': p[0] = (p[1],[])
            if p[0] not in self.P: self.P.append(p[0])
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
        self._normalize()
    def __repr__(self):
        out = ''
        for p in self.P:
            lhs = ' '.join([ self.symbols[i] if type(self.symbols[i]) == str
                            else str(self.symbols[i]) for i in p[0] ])
            rhs = ' '.join([ self.symbols[i] if type(self.symbols[i]) == str
                            else str(self.symbols[i]) for i in p[1] ])
            if not rhs: rhs = 'ε'
            out += lhs + ' = ' + rhs + '\n'
        return out
    # The wiki page on 'Kuroda normal form' mentions that any unrestricted language
    # is weakly equivalent to one where all rules are of the form:
    #   1.  AB = CD
    #   2.  A = BC
    #   3.  A = a
    #   4.  A = ε
    # If rule type 4. is ommitted, we have a normal form for context-sensitive
    # languages, whereas if only rule 1. is ommitted, we gain a normal form for
    # context-free languages.
    # To transform into normal form (note: I make a distinction between 'normal'
    # and 'canonical': canonical form being a unique, deterministic description
    # of a language, meaning a one-to-one mapping between language <-> canonical
    # grammar) follow algorithm on page 723 of 'Automata and Languages' by
    # Alexander Meduna:
    #   1.  In all productions of P, replace each occurrence of a terminal,
    #       a ∈ T, with a new non-terminal 'a and add 'a to N. Include
    #           'a = a
    #       into P.
    #   2.  Replace every production of the form
    #           A_1 ... A_m = B_1 ... B_n
    #       where n and m satisfy 0 <= n < m and 1 < m, with the following 2
    #       productions
    #           A_1 ... A_m = B_1 ... B_n {C}^(m-n)
    #           C = ε
    #       and add C into N.
    #   3.  Replace every production of the form
    #           A = B
    #       with the following 2 productions
    #           A = B C
    #           C = ε
    #       and add C into N.
    #   4.  Replace every production of the form
    #           A = B_1 ... B_n
    #       where 3 <= n, with the following n-1 productions
    #           A = B_1 B_2n
    #           B_2n = B_2 B_3n
    #               ...
    #           B_(n-2)n = B_(n-2) B_(n-1)n
    #           B_(n-1)n = B_(n-1) B_n
    #       and add the n-2 non-terminals
    #           B_2n ... B_(n-1)n
    #       into N.
    #   5.  Replace every production of the form
    #           A_1 ... A_m = B_1 ... B_n
    #       where 2 <= m and 3 <= n, with the following 2 productions
    #           A_1 A_2 = B_1 C
    #           C A_3 ... A_m = B_2 ... B_n
    #       and add C into N.  Repeat step 5 until all productions are in
    #       normal form.
    def _normalize(self):
        # 1.
        for i, s in enumerate(self.symbols[:]):
            if type(s) == int:
                self.symbols.append(s)
                self.symbols[i] = "'" + str(s)
                self.P.append(([i], [len(self.symbols)-1]))
        def new_symbol():
            i = len(self.symbols)
            self.symbols.append('_' + str(i))
            return i
        empty_nonterminal = new_symbol()
        self.P.append(([empty_nonterminal], []))
        for p in self.P[:]:
            # 2.
            if len(p[0]) > 1 and len(p[0]) > len(p[1]):
                p[1].extend([empty_nonterminal for _ in range(len(p[0]) - len(p[1]))])
            # 3.
            elif len(p[0]) == 1 and len(p[1]) == 1 \
                and type(self.symbols[p[0][0]]) == str \
                and type(self.symbols[p[1][0]]) == str:
                p[1].append(empty_nonterminal)
            # 4.
            elif len(p[0]) == 1 and len(p[1]) > 2:
                s = new_symbol()
                tail = p[1][1:]
                p[1][1:] = [s]
                while len(tail) > 2:
                    s_new = new_symbol()
                    self.P.append(([s], [tail[0], s_new]))
                    tail = tail[1:]
                    s = s_new
                else:
                    self.P.append(([s], tail))
            # 5.
            while len(p[0]) > 1 and len(p[1]) > 2:
                s = new_symbol()
                taill = p[0][2:]
                tailr = p[1][1:]
                del p[0][2:]
                p[1][1:] = [s]
                self.P.append(([s] + taill, tailr))
                p = self.P[-1]

