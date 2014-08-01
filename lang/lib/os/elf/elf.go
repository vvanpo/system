// Executable and Linkable Format
// http://refspecs.linux-foundation.org/elf/gabi4+/contents.html
package elf

import (
	"github.com/vvanpo/system/lang"
	"debug/elf"
	"io"
)

func Read(r io.ReaderAt) (format lang.file, err error) {
	f, err := elf.NewFile(r)
	// To return lang code format
	return
}

func Write(f lang.file) (w io.ReaderAt) {
	return
}
