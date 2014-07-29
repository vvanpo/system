
import asmlang, functools, re
from uuid import uuid1

class fd(bytearray):			# file descriptor
	def __init__(self, uuid):
		super().__init__(self)
		self.uuid = uuid
	def __hash__(self):
		return self.uuid

class kernel(object):
	timeslice = 100				# instructions per timeslice
	def __init__(self, init_process):
		self.fd = {fd(1)}		# set of files, 1 is stdout
		self.process = set()	# set of processes
		self._sched_proc(init_process)
		self._start()
	def _exec_instr(self, proc, instr):
		stmt = instr.pop(0)
		if stmt == "open":
			self._open(proc, instr)
		if stmt == "close":
			self._close(proc, instr)
		if stmt == "push":
			self._push(proc, instr)
		if stmt == "pop":
			pass
		if stmt == "copy":
			pass
		if stmt == "ifzero":
			pass
	def _open(self, proc, args):
		if re.match(r"^[a-z]+$", args[0]):		# BUG:  matches numbers like 'af'
			name = args.pop(0)
			if name == "stdout":
				if args and args[0] != 1:
					raise Exception("Incorrect uuid for special file stdout")
				uuid = "1"
		if args and re.match(r"^[0-9a-zA-Z]+$", args[0]):
			uuid = args[0]
			if not name:
				name = uuid
		elif not uuid:
			if not name:
				raise Exception("No name or uuid for segment creation")
			uuid = str(uuid1().int)
		for f in self.fd:
			if f.uuid == int(uuid):
				proc.open(f, name)
				return
		f = fd(int(uuid))
		self.fd.add(f)
		proc.open(f, name)
	def _close(self, proc, args):
		if args and re.match(r"^[a-z]+$", args[0]):
			name = args[0]
		elif args and re.match(r"^[0-9a-zA-Z]+$", args[0]):
			name = args[0]
		else:
			raise Exception("Incorrect close argument")
		del proc.segment[name]
	def _push(self, proc, args):
		proc.sp.value += 1

	def _sched_proc(self, code):
		p = process(code)
		self.process.add(p)
	def _start(self):
		while self.process:
			p = self.process.pop()
			for i in range(kernel.timeslice):
				instr = p.next()
				if instr == None: return
				self._exec_instr(p, instr)
			self.process.add(p)

class process(object):
	word_size = 8
	def __init__(self, code):
		self.instruction = asmlang.parse(code)
		self.segment = {"code": bytearray(len(self.instruction)),
						"main": bytearray()}
		self.ip = register("code")
		self.sp = register("main")
		self.fp = register("main")
	def open(self, fd, name):
		if name in self.segment:
			raise Exception("Segment name '" + name + "' already in use.")
		self.segment[name] = fd
	def next(self):
		if len(self.instruction) > self.ip.value:
			self.ip.value += 1
			return self.instruction[self.ip.value - 1]

class register(object):
	def __init__(self, segment, value=0):
		self.segment = segment
		self.value = value

