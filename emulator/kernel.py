
import asmlang, re

class fd(bytearray):			# file descriptor
	def __init__(self, fid):
		super().__init__(self)
		self.fid = fid
	def __hash__(self):
		return self.fid

class kernel(object):
	timeslice = 100				# instructions per timeslice
	def __init__(self, init_process):
		self.fd = {fd(1)}		# set of files, 1 is stdout
		self.process = set()	# set of processes
		self._sched_proc(init_process)
		self._start()
	def _open(self, proc, args):
		if re.match(r"^[a-z]+$", args[0]):		# BUG:  matches numbers like 'af'
			name = args.pop(0)
			if name == "stdout":
				if args and args[0] != 1:
					raise Exception("Incorrect uuid for special file stdout")
				fid = "1"
		if args and re.match(r"^[0-9a-zA-Z]+$", args[0]):
			fid = args[0]
			if not name:
				name = fid
		elif not fid:
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
	def __init__(self, code):
		self.instruction = asmlang.parse(code)
		self.segment = {2: bytearray(len(self.instruction)),
						3: bytearray()}
		self.ip = register(2)
		self.sp = register(3)
		self.fp = register(3)
	def open(self, fd):
		if fid in self.segment:
			raise Exception("Segment id " + str(fid) + " already in use.")
		self.segment[name] = fd
	def next(self):
		if len(self.instruction) > self.ip.value:
			self.ip.value += 1
			return self.instruction[self.ip.value - 1]

class register(object):
	def __init__(self, segment, value=0):
		self.segment = segment
		self.value = value

class machine(object):
	word_size = 8
	def __init__(self):
		self.kernel_mode = True
		self.ip = register(None)
		self.sp = register(None)
		self.fp = register(None)
	def exec(self, proc):
		pass
