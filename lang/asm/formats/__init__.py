import re, importlib, os

class section:
    def __init__(self, file, name, architecture, option_string):
        if not re.match(r"[\.\w-]+$", name): raise Exception("Invalid section name: " + name)
        # format instance
        self.file = file
        self.name = name
        self.architecture = architecture
        # label_name -> index_into_statements_list
        self.labels = {}
        self.parse_options(option_string)
    def __eq__(self, other):
        if self.name == other.name: return True
        return False
    def __repr__(self):
        return self.name + ':\n  ' + str(self.labels) + '\n  ' + repr(self.architecture)
    def add_label(self, name):
        # Labels need to be distinguishable from numbers, which should be defined
        # by architectures to have the form [0-9]+ or [0-9a-f]+h
        m = re.match(r'([0-9a-f]+h|[0-9]+)|([\w-]+)$', name)
        if not m or m.group(1):
            raise Exception("Invalid label: " + name)
        if name in self.labels:
            raise Exception("Duplicate label: " + self.name + '.' + name)
        self.labels[name] = len(self.architecture.statements)
    def add_statement(self, string): self.architecture.add_statement(string)

class bin_format:
    names = {}
    @classmethod
    def register(cls, name): cls.names[name] = cls
    @classmethod
    def new(cls, name): 
        if name not in cls.names:
            raise Exception("Invalid file format: " + name)
        return cls.names[name]()
    def __init__(self): self.sections = []
    def __repr__(self): return ''.join([ repr(s) for s in self.sections ])
    def add_section(self, section):
        if section in self.sections:
            raise Exception("Duplicate section name: " + str(section))
        self.sections.append(section)

for m in [ '.' + p[:-3] for p in os.listdir(*__path__) if p[0] != '.' and p[0] != '_' and p[-3:] == ".py" ]:
    importlib.import_module(m, __name__)

