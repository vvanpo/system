import asm

class instruction:
    @classmethod
    def from_string(cls, string):
        return string

# Register the 'ir' architecture name with ..asm package
asm.architecture.register('ir', instruction)

