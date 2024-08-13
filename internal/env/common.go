package blink

const ARCH_BIT = 32 << (^uint(0) >> 63)
const ARCH_BYTE = 8 * ARCH_BIT
const ARCH_WORD = 4 * ARCH_BIT
const ARCH_DWORD = 2 * ARCH_BIT
const ARCH_QWORD = 1 * ARCH_BIT
const ARCH_MAX = 64 * ARCH_BIT
const ARCH_MIN = 8 * ARCH_BIT

func Is64Bit() bool {
	return ARCH_BIT == 64
}

func IsRelease() bool {
	return _isRelease
}
