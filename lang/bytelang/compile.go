package main

func putWord(word uint) string {
	s := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		s[i] = byte(0xff | word)
		word >>= 8
	}
	return string(s)
}
