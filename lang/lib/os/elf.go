package os

const (
	// Type
	et_none   = 0
	et_rel    = 1
	et_exec   = 2
	et_dyn    = 3
	et_core   = 4
	et_loos   = 0xfe00
	et_hios   = 0xfeff
	et_loproc = 0xff00
	et_hiproc = 0xffff
	// Machine
	em_none   = 0
	em_386    = 3
	em_ppc64  = 21
	em_arm    = 40
	em_x86_64 = 62
	// Version
	ev_none    = 0
	ev_current = 1
	// Identification indices
	ei_mag0       = 0
	ei_mag1       = 1
	ei_mag2       = 2
	ei_mag3       = 3
	ei_class      = 4
	ei_data       = 5
	ei_version    = 6 // Value should be ev_current
	ei_osabi      = 7
	ei_abiversion = 8
	ei_pad        = 9
	ei_nident     = 16
	// Identification values
	ELFMAG0        = 0x7f
	ELFMAG1        = 'E'
	ELFMAG2        = 'L'
	ELFMAG3        = 'F'
	ELFCLASSNONE   = 0
	ELFCLASS32     = 1
	ELFCLASS64     = 2
	ELFDATANONE    = 0
	ELFDATA2LSB    = 1 // Little endian: MSB at lowest address
	ELFDATA2MSB    = 2 // Big endian: MSB at highest address
	ELFOSABI_NONE  = 0
	ELFOSABI_LINUX = 3
	// Special section indices
	shn_undef = 0 // Undefined, missing, or irrelevant section reference
	// If #-of-sections is greater than loreserve, elfHeader.shnum == shn_undef
	// and the actual number of entries is stored in secHeader.size of the first
	// entry
	shn_loreserve = 0xff00
	shn_loproc    = 0xff00
	shn_hiproc    = 0xff1f
	shn_loos      = 0xff20
	shn_hios      = 0xff3f
	shn_abs       = 0xfff1
	shn_common    = 0xfff2
	shn_xindex    = 0xffff
	shn_hireserve = 0xffff
	// Section types
	sht_null          = 0
	sht_progbits      = 1
	sht_symtab        = 2
	sht_strtab        = 3
	sht_rela          = 4
	sht_hash          = 5
	sht_dynamic       = 6
	sht_note          = 7
	sht_nobits        = 8
	sht_rel           = 9
	sht_shlib         = 10
	sht_dynsym        = 11
	sht_init_array    = 12
	sht_fini_array    = 13
	sht_preinit_array = 14
	sht_group         = 17
	sht_symtab_shndx  = 18
	sht_loos          = 0x60000000
	sht_hios          = 0x6fffffff
	sht_loproc        = 0x70000000
	sht_hiproc        = 0x7fffffff
	sht_louser        = 0x80000000
	sht_hiuser        = 0xffffffff
	// Section attribute flags
	shf_write            = 0x1
	shf_alloc            = 0x2
	shf_execinstr        = 0x4
	shf_merge            = 0x10
	shf_strings          = 0x20
	shf_info_link        = 0x40
	shf_link_order       = 0x80
	shf_os_nonconforming = 0x100
	shf_group            = 0x200
	shf_tls              = 0x400
	shf_maskos           = 0x0ff00000
	shf_maskproc         = 0xf0000000
)

// Representation of an ELF object file
// http://refspecs.linux-foundation.org/elf/gabi4+/contents.html
type Elf struct {
	ident     [ei_nident]byte // Identification
	elfHeader interface{}
}

func New(class, data, osabi, abiversion byte) (e *Elf) {
	e = new(Elf)
	e.ident[ei_mag0] = ELFMAG0
	e.ident[ei_mag1] = ELFMAG1
	e.ident[ei_mag2] = ELFMAG2
	e.ident[ei_mag3] = ELFMAG3
	e.ident[ei_class] = class
	e.ident[ei_data] = data
	e.ident[ei_version] = ev_current
	e.ident[ei_osabi] = osabi
	if osabi == ELFOSABI_NONE {
		e.ident[ei_abiversion] = 0
	} else {
		e.ident[ei_abiversion] = abiversion
	}
	switch class {
	case ELFCLASS64:
		e.elfHeader = elfHeader64{}
	case ELFCLASS32:
		e.elfHeader = elfHeader32{}
	}
	return
}

func (e *Elf) Compose() (s string) {
	return
}

type elfHeader32 struct {
	typ       half32 // Type
	machine   half32 // Machine
	version   word32 // Version
	entry     addr32 // Virtual address entry to begin process (0 == no entry point)
	phoff     off32  // Program header table's offset in bytes	(0 == no program header table)
	shoff     off32  // Section header table's offset in bytes	(0 == no section header table)
	flags     word32 // Processor-specific flags
	ehsize    half32 // ELF header size in bytes
	phentsize half32 // Program header table entry size in bytes
	phnum     half32 // Number of entries in program header table
	shentsize half32 // Section header (one entry in section header table) size in bytes
	shnum     half32 // Number of entries in section header table
	shstrndx  half32 // Section header table index of section name string table entry
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
