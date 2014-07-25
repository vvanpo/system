
class fd(list):				# file descriptor
	def __init__(self, uuid)
		self.uuid = uuid
		list.__init__(self)

class kernel(object):
	def __init__(self, init_process):
		self.fd = {}		# set of files
	def open(self, uuid):
		for f in self.fd:
			if f.uuid == uuid:
				return f
		f = fd(uuid)
		self.fd.add(f)
		return f

