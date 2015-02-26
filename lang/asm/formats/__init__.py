import re, importlib, os

class section:
    def __init__(self, file, name, architecture, option_string):
        if not re.match(r"[\.\w-]+$", name): raise Exception("Invalid section name: " + name)
        self.file = file
        self.name = name
        self.architecture = architecture
        self.statements = []
        self.parse_options(option_string)
    def __eq__(self, other):
        if self.name == other.name: return True
        return False
    def __str__(self):
        return self.name
    def add_statement(self, statement):
        self.statements.append(self.architecture.statement(statement))

class bin_format:
    names = {}
    @classmethod
    def register(cls, name):
        cls.names[name] = cls
    @classmethod
    def new(cls, name): 
        if name not in cls.names:
            raise Exception("Invalid file format: " + name)
        return cls.names[name]()
    def __init__(self):
        self.sections = []
    def add_section(self, name):
        if name in self.sections: raise Exception("Duplicate section name: " + str(s))
        self.sections.append(name)

for m in [ '.' + p[:-3] for p in os.listdir(*__path__) if p[0] != '.' and p[0] != '_' and p[-3:] == ".py" ]:
    importlib.import_module(m, __name__)

