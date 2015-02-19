import re

class instruction:
    def __init__(self, opcode, mnemonic, operands, encoding, compat, long):
        self.mnemonic = mnemonic

def load_instructions():
    instructions = set()
    with open("instructions.txt") as f:
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
                    instructions.add(instruction(opcode, mnemonic, operands, enc, compat, long))
            else:
                mnemonic = line.strip()
                encoding.clear()
                sec_opcodes = False
    return instructions

instructions = load_instructions()

for i in instructions:
    print(i.mnemonic)
