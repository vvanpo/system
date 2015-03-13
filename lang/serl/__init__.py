import re, os

# serl: Serialization Language
# A data format that can map to sets, dicts, lists, etc. and is human-readable.
# The character set, for now, UTF-8 by default.
#   file = list? end-of-file
#   end-of-file = "——" newline
#   list = plain-list | dashed-list
#   plain-list = (item newline)+
#   dashed-list = ("- " item newline)+
#   item = anchor? pair | scalar | list
#   pair = simple-pair | complex-pair
#   simple-pair = scalar ": " anchor? scalar | list
#   complex-pair = "?" newline | " " anchor? list ":" anchor? scalar | list
#   
class file:
    @classmethod
    def from_string(cls, string):
        pass
    @classmethod
    def from_variable(cls, var):
        pass
    def __repr__(self):
        return str(self.root)

def parse(source):
    return file.from_string(source)

def dump(var):
    return file.from_variable(var)
