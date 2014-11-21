
.text

	.align 4
	.long 0x1BADB002
	.long 0
	.long -0x1BADB002

.global entry
.extern boot

entry:
	call boot
