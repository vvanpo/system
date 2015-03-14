
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
# whitespace (can cross line boundaries), and terminals are indicated with double
# quotes (although a character set must be defined beforehand).
#   start = production+
#   production = string nonterminal string "=" string
#   string = symbol* | empty-string
#     empty-string = "ε"
# When a grammar is processed, non-terminals are either can be either "capturing"
# if they begin with a letter, or "non-capturing" if they begin with an underscore.
# Capturing non-terminals will be used by the generated parser to return matching
# strings.  Non-capturing non-terminals will not be available for use and may be
# transformed to different symbol(s).  The grammar class transforms the input into
# a canonical form, and determines which category in the hierarchy it falls into.

class grammar:
    pass
