from asm.formats import section, bin_format

class section(section):
    def parse_options(self, string):
        self.align = 0
        self.start = None   # section start-address is determined by the length
                            # of previous sections
        if not string: return
        option = string.partition('=')
        if option[0] == "start": self.start = int(option[2])
        elif option[0] == "align": self.align = int(option[2])
        else: raise Exception("Invalid section option '" + option[0] + "' for bin format")

class bin(bin_format):
    def new_section(self, *args):
        s = section(self, *args)
        self.add_section(s)
        return s
    def calculate_addr(self):
        for s in self.sections:
            pass

bin.register('bin')

