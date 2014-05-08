package main

type symbol string

type symTable map[symbol]struct {
	symType
	ref uint64
}

type symType struct {
	size uint
}
