from asm.formats import section, bin_format

class section(section):
    def parse_options(self, string):
        self.align = 0
        self.start = None   # section start-address is determined by the length
                            # of previous sections
        if not string: return
        m = re.match(r'(?:\s*(?:align=([0-9]+)|start=([0-9]+)))+', string)
        if m and m.group(1): self.align = int(m.group(1))
        if m and m.group(2): self.start = int(m.group(2))
        if not m and string.strip():
            raise Exception("Invalid section option '" + string.strip() + "' for bin format")

class bin(bin_format):
    def add_section(self, *args):
        s = section(self, *args)
        super().add_section(s)
        return s
    def assemble(self):
        for i in range(len(self.sections)):
            s = self.sections[i]

bin.register('bin')

