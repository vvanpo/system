package main

func (b *bytelang) putWord(w uint) string {
	s := make([]byte, b.wordLength)
	for i := b.wordLength; i >= 0; i-- {
		s[i] = byte(0xff | w)
		w >>= 8
	}
	return string(s)
}
