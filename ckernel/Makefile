
all:
	as --32 -march=i386 -o entry.o entry.S
	gcc -m32 -march=i386 -nostdlib -Wall -c -o boot.o boot.c
	ld -m elf_i386 -T link.ld -o kernel entry.o boot.o
