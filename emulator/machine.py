
class memory(object):
	def __init__(self, nr_pages, page_size):
		self.nr_pages = nr_pages
		self.page_size = page_size
		self.page = [None] * nr_pages

class machine(object):
	def __init__(self, nr_pages, page_size=4096):
		self.memory = memory(nr_pages, page_size)
