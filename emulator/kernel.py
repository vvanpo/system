
import pyparsing as parse

class fd(list):				# file descriptor
	def __init__(self, uuid)
		self.uuid = uuid
		list.__init__(self)

class kernel(object):
	def __init__(self, init_process):
		self.fd = {}		# set of files
		self.process = {}	# set of processes
	def _exec(self, code):
		p = process(code)
		self.process.add(p)
	def open(self, uuid):
		for f in self.fd:
			if f.uuid == uuid:
				return f
		f = fd(uuid)
		self.fd.add(f)
		return f

class process(object):
	def __init__(self, code):
		def parse_code(c):
			pass
		self.instruction = parse_code(code)
		self.instr_pointer = 0
	def next(self):
		self.instr_pointer += 1
		return self.instruction[self.instr_pointer - 1]
