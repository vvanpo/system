import re, os.path

class instruction:
    def __init__(self, opcode, operands, encoding, compat, long):
        pass

# Assembles instruction into binary format
def assemble(source):
    # Call nasm for now
    pass
# Disassembles binary instruction
def disassemble(binary):
    # Call ndisasm for now
    pass

def load_instructions():
    instructions = {}
    dir = os.path.dirname(__file__)
    with open(dir + "/instructions.txt") as f:
        mnemonic = ""
        encoding = {}
        sec_opcodes = False
        for line in f:
            if line[0] == "#": continue
            if line[:2] == "  ":
                if line[2:5] == "---":
                    sec_opcodes = True
                elif not sec_opcodes:
                    l = line.split()
                    encoding[l[0]] = tuple(l[1:])
                else:
                    opcode = line[2:28].strip()
                    operands = line[28:48].strip()
                    enc = encoding[line[48:52].strip()]
                    compat = line[52:56].strip()
                    long = line[56:].strip()
                    instructions[mnemonic].add(instruction(opcode, operands, enc, compat, long))
            else:
                mnemonic = line.strip()
                instructions[mnemonic] = set()
                encoding.clear()
                sec_opcodes = False
    return instructions

instructions = load_instructions()

def mov(width, source, dest):
    print(instructions["MOV"])
