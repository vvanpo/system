
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
#   S = production+
#     production = _string? _space non-terminal _space _string? "=" \
#                   _string | empty-string _newline
#       _string = _space? group | concat _space?
#         group = ("(" _space? group | concat _space? ")") | _symbol
#                   | (group (zero-or-one | zero-or-more | one-or-more | _range))
#           _symbol = terminal | non-terminal
#           zero-or-one = "?"
#           zero-or-more = "*"
#           one-or-more = "+"
#           _range = "{" range-exact | (range-low "," range-high?) | ("," range-high) "}"
#             range-exact = _digit+
#             range-low = _digit+
#             range-high = _digit+
#         concat = group | option (_space group | option)+
#           option = group (_space? "|" _space? group)+
#       _space = " "+
#         ###_escaped-newline = "\" _newline
#       non-terminal = "_"? (_letter | _digit) (_letter | _digit | "-")*
#       terminal = ('"' _character '"') | _hex
#         _hex = ("0" .. "9" | "a" .. "f")+
#       empty-string = _space? "ε" _space?
# When a grammar is processed, non-terminals are either can be either "capturing"
# if they begin with a letter, or "non-capturing" if they begin with an underscore.
# Capturing non-terminals will be used by the generated parser to return matching
# strings.  Non-capturing non-terminals will not be available for use and may be
# transformed to different symbol(s).  The grammar class transforms the input into
# a normal form (with the option to dump this form as either simple productions or
# using the notation above), and determines which category in the hierarchy it falls
# into.

class grammar:
    def __init__(self, source, character_set='utf-8'):
        pass
    def dump(self, simple=False):
        pass

