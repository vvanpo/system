import asm

class instruction:
    models = {}
    mnemonics = {}
    @classmethod
    def _load(cls, filepath):
        with open(filepath) as f:
            mnemonic = ""
            encoding = {}
            opcode_section = False
            for line in f:
                if line[0] == "#": continue
                if line[:2] == "  ":
                    if line[2:5] == "---":
                        opcode_section = True
                    elif not opcode_section:
                        l = line.split()
                        encoding[l[0]] = tuple(l[1:])
                    else:
                        opcode = line[2:28].strip()
                        operands = line[28:48].strip()
                        enc = encoding[line[48:52].strip()]
                        compat = line[52:56].strip()
                        long = line[56:60].strip()
                        instruction = (opcode, operands, enc, compat, long)
                        cls.mnemonics[mnemonic].add(instruction)
                        if len(line) > 60:
                            intro = line[60:].strip()
                            if intro in cls.models:
                                cls.models[intro].add(instruction)
                            else:
                                cls.models[intro] = set((instruction,))
                else:
                    mnemonic = line.strip()
                    cls.mnemonics[mnemonic] = set()
                    encoding.clear()
                    opcode_section = False
    @classmethod
    def from_string(cls, string):
        pass

instruction._load(__path__[0] + "/instructions.txt")
print(instruction.models)

# Register with architecture from ..asm
for m in instruction.models.keys():
    asm.architecture.register(m, instruction)

