
import machine

class process(object):
	def __init__(self):
		self.page_table = {}

class stream(object):
	def __init__(self)
		self.size = 0
		self.index = 0

class kernel(process):
	def __init__(self, machine):
		self.machine = machine
		self.process = {}	# uuid: process
		self.stream = {}	# stream: [page_nrs]
	def read(self, stream):
		pass
	def write(self, stream):
		pass
	def load(self, file):
		pass

