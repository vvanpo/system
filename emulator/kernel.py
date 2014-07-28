
import asmlang

class fd(list):					# file descriptor
	def __init__(self, uuid):
		self.uuid = uuid
		list.__init__(self)

class kernel(object):
	def __init__(self, init_process):
		self.fd = set()			# set of files
		self.process = set()	# set of processes
		self._exec(init_process)
	def _exec(self, code):
		p = process(code)
		self.process.add(p)
	def open(self, uuid=None):
		for f in self.fd:
			if f.uuid == uuid:
				return f
		f = fd(uuid)
		self.fd.add(f)
		return f

class process(object):
	def __init__(self, code):
		self.instruction = asmlang.parse(code)
		self.instr_ptr = 0
		self.stack_ptr = 0
		self.frame_ptr = 0
	def next(self):
		self.instr_ptr += 1
		return self.instruction[self.instr_ptr - 1]
