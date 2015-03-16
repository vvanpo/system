import re

# An order map has all the same operations as a regular list, but items that are
# pairs can be indexed by their key (which returns their value), or their index
# (which returns the key-value pair)
class map(list):
    def __init__(self, *args, unique_keys=False, ordered=True):
        super().__init__(*args)
        self.unique_keys = unique_keys
        self.ordered = ordered
    def __getitem__(self, key):
        pass

# serl: Serialization Language
# A data format that can map to sets, dicts, lists, etc. and is human-readable.
# The character set, for now, UTF-8 by default.
# For now, indent tokens are handled outside the grammar, until I can figure out
# how to write a proper context-sensitive grammar, and corresponding parser.
#   S = map? end-of-file
#     end-of-file = "——" newline
#     map = plain-map | dashed-map
#       plain-map = (item newline)+
#       dashed-map = (("- " item) | ("? " key (newline INDENT ":" value)?) newline)+
#   
class file:
    def __init__(self):
        self.root = map()
    def __repr__(self):
        return str(self.root)
    @classmethod
    def from_string(cls, string):
        self = cls()
        return self
    @classmethod
    def from_variable(cls, var):
        pass

def parse(source):
    return file.from_string(source)

def dump(var):
    return file.from_variable(var)

