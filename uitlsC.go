package blink

import (
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

type (
	// char
	Char byte
	// singed char
	SChar int8
	// unsigned char
	UChar uint8
	// short
	Short int16
	// unsigned short
	UShort uint16
	// int
	Int int32
	// unsigned int
	UInt uint32
	// long
	Long int32
	// unsigned long
	ULong uint32
	// long int
	LongInt int64
	// unsigned long int
	ULongInt uint64
	// long long
	LongLong int64
	// unsigned long long
	ULongLong uint64
	// float
	Float float32
	// double
	Double float64
	// long double
	LongDouble float64
	// wchar_t
	Wchar_t uint16
	// utf8
	Utf8 byte
)

// alias
type (
	// void *
	VoidPtr = uintptr
	AnyPtr  = uintptr
)

func StringToPtr(s string) uintptr {
	p, err := windows.BytePtrFromString(s)
	if err != nil {
		return 0
	}

	return (uintptr)(unsafe.Pointer(p))
}

func PtrToString(p uintptr) string {
	return windows.BytePtrToString(AssertType[uint8](p))
}

func StringToChar(s string) []Char {
	bytes, err := windows.ByteSliceFromString(s)
	if err != nil {
		return nil
	}
	return *((*[]Char)(unsafe.Pointer(&bytes)))
}

func StringToWCharPtr(s string) uintptr {
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		return 0
	}
	return uintptr(unsafe.Pointer(p))
}

func StringToWcharU16Ptr(s string) *uint16 {
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		*p = 0
	}
	return p
}

func StringToU16Arr(s string) []uint16 {
	return utf16.Encode([]rune(s))
}

func PtrWCharToString(p uintptr) string {
	return windows.UTF16PtrToString(AssertType[uint16](p))
}

// callback = func(args ...uintptr) uintptr
func CallbackToPtr(callback any) uintptr {
	return syscall.NewCallbackCDecl(callback)
}

func BoolToPtr(b bool) uintptr {

	if b {
		return 1
	}

	return 0
}

func PtrToBool(p uintptr) bool {
	return p != 0
}

func AssertType[T any](ptr uintptr) *T {
	return (*T)(unsafe.Pointer(ptr))
}

type KnownCType interface {
	// char
	Char |
		// singed char
		SChar |
		// unsigned char
		UChar |
		// short
		Short |
		// unsigned short
		UShort |
		// int
		Int |
		// unsigned int
		UInt |
		// long
		Long |
		// unsigned long
		ULong |
		// long int
		LongInt |
		// unsigned long int
		ULongInt |
		// long long
		LongLong |
		// unsigned long long
		ULongLong |
		// float
		Float |
		// double
		Double |
		// long double
		LongDouble |
		// wchar_t
		Wchar_t |
		// utf8
		Utf8
}

func LOWORD(dwValue uint32) uint16 {
	return uint16(dwValue & 0xFFFF)
}

func CopyString(src uintptr, n int) string {
	return string(CopyBytes(src, n))
}

func CopyBytes(src uintptr, n int) []byte {
	if n == 0 {
		return make([]byte, 0)
	}

	byts := make([]uint8, n)
	for i := 0; i < n; i++ {
		byts[i] = *(*uint8)(unsafe.Pointer(src + uintptr(i)))
	}

	return byts
}
