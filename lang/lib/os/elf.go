package os

const (
	// Type
	t_none   int = 0
	t_rel        = 1
	t_exec       = 2
	t_dyn        = 3
	t_core       = 4
	t_loos       = 0xfe00
	t_hios       = 0xfeff
	t_loproc     = 0xff00
	t_hiproc     = 0xffff
	// Machine
	m_none   = 0
	m_386    = 3
	m_ppc64  = 21
	m_arm    = 40
	m_x86_64 = 62
	// Version
	v_none    = 0
	v_current = 1
	// Identification indices
	i_mag0       = 0
	i_mag1       = 1
	i_mag2       = 2
	i_mag3       = 3
	i_class      = 4
	i_data       = 5
	i_version    = 6 // Value should be ev_current
	i_osabi      = 7
	i_abiversion = 8
	i_pad        = 9
	i_nident     = 16
	// Identification values
	mag0        = 0x7f
	mag1        = 'E'
	mag2        = 'L'
	mag3        = 'F'
	classnone   = 0
	class32     = 1
	class64     = 2
	datanone    = 0
	data2lsb    = 1 // Little endian: MSB at lowest address
	data2msb    = 2 // Big endian: MSB at highest address
	osabi_none  = 0
	osabi_linux = 3
)

// Representation of an ELF object file
// http://refspecs.linux-foundation.org/elf/gabi4+/contents.html
type elf struct{}

type elfHeader interface {
	ident(i [i_nident]byte) // Identification
	typ()                   // Type
	machine()               // Machine
	version()               // Version
	entry()                 // Virtual address entry to begin process (0 == no entry point)
	phoff()                 // Program header table's offset in bytes	(0 == no program header table)
	shoff()                 // Section header table's offset in bytes	(0 == no section header table)
	flags()                 // Processor-specific flags
	ehsize()                // ELF header size in bytes
	phentsize()             // Program header table entry size in bytes
	phnum()                 // Number of entries in program header table
	shentsize()             // Section header (one entry in section header table) size in bytes
	shnum()                 // Number of entries in section header table
	shstrndx()              // Section header table index of section name string table entry
}

type elfHeader32 struct {
	ident     [i_nident]byte
	typ       half32
	machine   half32
	version   word32
	entry     addr32
	phoff     off32
	shoff     off32
	flags     word32
	ehsize    half32
	phentsize half32
	phnum     half32
	shentsize half32
	shnum     half32
	shstrndx  half32
}

type secHeader32 struct {
	name      word32
	typ       word32
	flags     word32
	addr      addr32
	offset    off32
	size      word32
	link      word32
	info      word32
	addralign word32
	entsize   word32
}

type elfHeader64 struct {
	ident     [i_nident]byte
	typ       half64
	machine   half64
	version   word64
	entry     addr64
	phoff     off64
	shoff     off64
	flags     word64
	ehsize    half64
	phentsize half64
	phnum     half64
	shentsize half64
	shnum     half64
	shstrndx  half64
}

type secHeader64 struct {
	name      word64
	typ       word64
	flags     xword64
	addr      addr64
	offset    off64
	size      xword64
	link      word64
	info      word64
	addralign xword64
	entsize   xword64
}

// 32-bit data types
type addr32 [4]byte
type off32 [4]byte
type half32 [2]byte
type word32 [4]byte
type sword32 [4]byte

// 64-bit data types
type addr64 [8]byte
type off64 [8]byte
type half64 [2]byte
type word64 [4]byte
type sword64 [4]byte
type xword64 [8]byte
type sxword64 [8]byte
