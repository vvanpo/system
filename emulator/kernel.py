
class process(object):
	def __init__(self):
		self.page_table = {}

class stream(object):
	def __init__(self)
		self.size = 0
		self.index = 0
		self.page = []		# pages, None means not loaded

class kernel(process):
	def __init__(self, machine):
		self.machine = machine
		self.process = {}	# uuid: process
		self.stream = {}	# uuid: stream
	def find_page(self):
		pass
	def read(self, stream, size=None):
		# data is of type bytes or bytearray
		# When size is None read() returns the entire buffer
		pass
	def write(self, stream, data):
		# data is of type bytes or bytearray
		page_sz = self.machine.memory.page_sz
		fragment = stream.size % page_sz
		nr_pages = -(-stream.size // page_sz)
		if stream.index > fragment:
			start = stream.index - fragment
			wraps = True
		else:
			start = stream.index + stream.size
		# start is the index of the end of the current stream
		# if start < stream.index, the stream wraps
		# fragment is the space left in the currently allocated pages
		if len(data) > fragment:
			nr_new_pages = -(-(len(data) - fragment) // page_sz)
			for i in range(nr_new_pages):
				stream.page.append(self.find_page())
			stream.page[nr_pages][0:start] = stream.page[0][0:start]
			start += nr_pages * page_sz - 1
		stream.size += len(data)
		
	def exec(self, fd):
		pass

