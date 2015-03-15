
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
#   start = production+
#     production = string? space non-terminal space string? "=" string | empty-string newline
#       string = space? group | concat space?
#         group = ("(" space? group | concat space? ")") | symbol
#                   | (group (zero-or-one | zero-or-more | one-or-more | range))
#           symbol = terminal | non-terminal
#           zero-or-one = "?"
#           zero-or-more = "*"
#           one-or-more = "+"
#           range = "{" digit | (digit+ "," digit*) | ("," digit+) "}"
#         concat = (group | option (space group | option)+)
#           option = (group (space "|" space group)+) | symbol
#       space = " "+
#         ###escaped-newline = "\" newline
#       non-terminal = "_"? (letter | digit) (letter | digit | "-")*
#       empty-string = space? "ε" space?
# When a grammar is processed, non-terminals are either can be either "capturing"
# if they begin with a letter, or "non-capturing" if they begin with an underscore.
# Capturing non-terminals will be used by the generated parser to return matching
# strings.  Non-capturing non-terminals will not be available for use and may be
# transformed to different symbol(s).  The grammar class transforms the input into
# a normal form (with the option to dump this form as either simple productions or
# using the notation above), and determines which category in the hierarchy it falls
# into.

class grammar:
    def __init__(self, source):
        pass
    def dump(self, simple=False):
        pass

