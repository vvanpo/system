
import asmlang, functools

class fd(bytearray):			# file descriptor
	def __init__(self, uuid):
		super().__init__(self)
		self.uuid = uuid

class kernel(object):
	def __init__(self, init_process):
		self.fd = set(fd(1))	# set of files, 1 is stdout
		self.process = set()	# set of processes
		self._exec_proc(init_process)
		self._start()
	def _exec_instr(self, proc, instr):
		stmt = instr.pop(0)
		if stmt == "open":
			pass
		if stmt == "close":
			pass
		if stmt == "push":
			pass
		if stmt == "pop":
			pass
		if stmt == "copy":
			pass
		if stmt == "ifzero":
			pass
	def _exec_proc(self, code):
		p = process(code)
		self.process.add(p)
	def _start(self):
		while len(self.process) > 0:
			p = self.process.pop()
			for i in range(100):	# 100 instructions per timeslice
				instr = p.next()
				if instr == None: return
				self._exec_instr(p, instr)
			self.process.add(p)
	def open(self, uuid=None):
		for f in self.fd:
			if f.uuid == uuid:
				return f
		f = fd(uuid)
		self.fd.add(f)
		return f

class process(object):
	word_size = 8
	def __init__(self, code):
		self.instruction = asmlang.parse(code)
		self.segment = {"code": bytearray(len(self.instruction)),
						"main": bytearray()}
		self.ip = register("code")
		self.sp = register("main")
		self.fp = register("main")
	def next(self):
		if len(self.instruction) > self.ip.value:
			self.ip.value += 1
			return self.instruction[self.ip.value - 1]

class register(object):
	def __init__(self, segment, value=0):
		self.segment = segment
		self.value = value

