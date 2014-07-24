
class memory(object):
	def __init__(self, nr_pages, page_sz):
		self.nr_pages = nr_pages
		self.page_sz = page_sz
		self.page = [None] * nr_pages	# pages are bytearrays

class machine(object):
	def __init__(self, nr_pages, page_sz=4096):
		self.memory = memory(nr_pages, page_sz)
