// Executable and Linkable Format
// http://refspecs.linux-foundation.org/elf/gabi4+/contents.html
package elf

import (
	"bytes"
	"debug/elf"
	"github.com/vvanpo/system/lang"
	"io"
)

func Read(r io.ReaderAt) (file *lang.File, err error) {
/*	f, err := elf.NewFile(r)
	if err != nil {
		return
	}
	d := f.fileHeader.Data
	c := f.fileHeader.Class
	m := f.fileHeader.Machine
*/	file = new(lang.File)
	// To return lang code format
	return
}

func Write(f lang.File, d elf.Data, c elf.Class, m elf.Machine) (r *bytes.Reader) {
	var b []byte
	b = append(b, elf.ELFMAG...)
	r = bytes.NewReader(b)
	return
}

