package blink

import (
	"syscall"
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

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 |
		Char | SChar | UChar | Short | UShort | Int | UInt | Long | ULong | LongInt | ULongInt | LongLong |
		ULongLong | Float | Double | LongDouble | Wchar_t | Utf8 |
		uintptr
}

func Read[T Numeric](ptr uintptr) []T {
	var bytes []T
	var byt T
	for {
		byt = *((*T)(unsafe.Pointer(ptr)))

		if byt == 0 {
			break
		}

		bytes = append(bytes, byt)
		ptr++
	}

	return bytes
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
