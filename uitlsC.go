package blink

import (
	"bytes"
	"encoding/binary"
	"errors"
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
func CallbackToPtr(callback interface{}) uintptr {
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

func AssertType[T interface{}](ptr uintptr) *T {
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

func HIWORD(dwValue uint32) uint16 {
	return uint16((dwValue >> 16) & 0xFFFF)
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

const ARCH_BIT = 32 << (^uint(0) >> 63)

type LangConfig struct {
	IsLittleEndian bool
}

func IntToBytes(value int) []byte {
	if ARCH_BIT == 64 {
		return Int64ToBytes(int64(value))
	}

	return Int32ToBytes(int32(value))
}

func Int16ToBytes(value int16) []byte {
	bytes := make([]byte, 2)
	*(*int16)(unsafe.Pointer(&bytes[0])) = value
	return bytes
}

func Int32ToBytes(value int32) []byte {
	bytes := make([]byte, 4)
	*(*int32)(unsafe.Pointer(&bytes[0])) = value
	return bytes
}

func Int64ToBytes(value int64) []byte {
	bytes := make([]byte, 8)
	*(*int64)(unsafe.Pointer(&bytes[0])) = value
	return bytes
}

func UintToBytes(value uint) []byte {
	if 64 == ARCH_BIT {
		return Uint64ToBytes(uint64(value))
	}

	return Uint32ToBytes(uint32(value))
}

func Uint16ToBytes(value uint16) []byte {
	bytes := make([]byte, 2)
	*(*uint16)(unsafe.Pointer(&bytes[0])) = value
	return bytes
}

func Uint32ToBytes(value uint32) []byte {
	bytes := make([]byte, 4)
	*(*uint32)(unsafe.Pointer(&bytes[0])) = value
	return bytes
}

func Uint64ToBytes(value uint64) []byte {
	bytes := make([]byte, 8)
	*(*uint64)(unsafe.Pointer(&bytes[0])) = value
	return bytes
}

func Float32ToBytes(value float32) []byte {
	bytes := make([]byte, 4)
	*(*float32)(unsafe.Pointer(&bytes[0])) = value
	return bytes
}

func Float64ToBytes(value float64) []byte {
	bytes := make([]byte, 8)
	*(*float64)(unsafe.Pointer(&bytes[0])) = value
	return bytes
}

func IntFromBytes(src []byte, startIndex int, isLittleEndian bool) (int, error) {
	if ARCH_BIT == 64 {
		return numericFromBytes[int](src, startIndex, 8, isLittleEndian)
	}
	return numericFromBytes[int](src, startIndex, 4, isLittleEndian)
}

func Int16FromBytes(src []byte, startIndex int, isLittleEndian bool) (int16, error) {
	return numericFromBytes[int16](src, startIndex, 2, isLittleEndian)
}

func Int32FromBytes(src []byte, startIndex int, isLittleEndian bool) (int32, error) {
	return numericFromBytes[int32](src, startIndex, 4, isLittleEndian)
}

func Int64FromBytes(src []byte, startIndex int, isLittleEndian bool) (int64, error) {
	return numericFromBytes[int64](src, startIndex, 8, isLittleEndian)
}

func UintFromBytes(src []byte, startIndex int, isLittleEndian bool) (uint, error) {
	if ARCH_BIT == 64 {
		return numericFromBytes[uint](src, startIndex, 8, isLittleEndian)
	}
	return numericFromBytes[uint](src, startIndex, 4, isLittleEndian)
}

func Uint16FromBytes(src []byte, startIndex int, isLittleEndian bool) (uint16, error) {
	return numericFromBytes[uint16](src, startIndex, 2, isLittleEndian)
}

func Uint32FromBytes(src []byte, startIndex int, isLittleEndian bool) (uint32, error) {
	return numericFromBytes[uint32](src, startIndex, 4, isLittleEndian)
}

func Uint64FromBytes(src []byte, startIndex int, isLittleEndian bool) (uint64, error) {
	return numericFromBytes[uint64](src, startIndex, 8, isLittleEndian)
}

func Float32FromBytes(src []byte, startIndex int, isLittleEndian bool) (float32, error) {
	return numericFromBytes[float32](src, startIndex, 4, isLittleEndian)
}

func Float64FromBytes(src []byte, startIndex int, isLittleEndian bool) (float64, error) {
	return numericFromBytes[float64](src, startIndex, 8, isLittleEndian)
}

type Numeric interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func numericFromBytes[T Numeric](src []byte, startIndex int, bit int, isLittleEndian bool) (T, error) {
	if src == nil {
		return 0, errors.New("值不能为空")
	}
	if startIndex >= len(src) {
		return 0, errors.New("数组越界")
	}
	if startIndex > len(src)-bit {
		return 0, errors.New("数组太小")
	}

	byt := src[startIndex : startIndex+bit]

	buf := bytes.NewBuffer(byt)

	var num T
	var err error

	if isLittleEndian {
		err = binary.Read(buf, binary.LittleEndian, &num)
	} else {
		err = binary.Read(buf, binary.BigEndian, &num)
	}

	if err != nil {
		return 0, err
	}

	return num, nil
}
