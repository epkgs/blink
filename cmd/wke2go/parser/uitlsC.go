package parser

import (
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

///////////////////////////////////////////////////////////////////////////////////////////////
// C 语言类型               CGO 类型         Go 语言类型
// char                    C.char           byte
// singed char             C.schar          int8
// unsigned char           C.uchar          uint8
// short                   C.short          int16
// unsigned short          C.ushort         uint16
// int                     C.int            int32
// unsigned int            C.uint           uint32
// long                    C.long           int32
// unsigned long           C.ulong          uint32
// long long int           C.longlong       int64
// unsigned long long int  C.ulonglong      uint64
// float                   C.float          float32
// double                  C.double         float64
// size_t                  C.size_t         uintptr
///////////////////////////////////////////////////////////////////////////////////////////////

type (
	Int8    int8
	Uint8   uint8
	Int16   int16
	Uint16  uint16
	Int32   int32
	Uint32  uint32
	Int64   int64
	Uint64  uint64
	Float32 float32
	Float64 float64
	Bool    bool
	Uintptr uintptr

	Byte = Uint8
)

type (
	// char
	Char Byte
	// singed char
	SChar Int8
	// unsigned char
	UChar Uint8
	// short
	Short Int16
	// unsigned short
	UShort Uint16
	// int
	Int Int32
	// unsigned int
	UInt Uint32
	// long
	Long Int32
	// unsigned long
	ULong Uint32
	// long long
	LongLong Int64
	// unsigned long long
	ULongLong Uint64
	// float
	Float Float32
	// double
	Double Float64
	// void*
	AnyPtr Uintptr
	// utf8
	String string
	// wchar_t
	WString string
)

func StringToPtr(s string) uintptr {
	p, err := windows.BytePtrFromString(s)
	if err != nil {
		return 0
	}

	return (uintptr)(unsafe.Pointer(p))
}

func PtrToString(p uintptr) string {
	return windows.BytePtrToString(AssertType[byte](p))
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

func CopyString(src uintptr, n int) string {
	return string(CopyBytes(src, n))
}

func CopyBytes(src uintptr, n int) []byte {
	if n == 0 {
		return make([]byte, 0)
	}

	byts := make([]byte, n)
	for i := 0; i < n; i++ {
		byts[i] = *(*byte)(unsafe.Pointer(src + uintptr(i)))
	}

	return byts
}

// //////////////////////////////////////////////////
// Int8
func (c *Int8) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Int8) FromPtr(p uintptr) Int8 {
	*c = Int8(p)
	return *c
}

// Uint8
func (c *Uint8) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Uint8) FromPtr(p uintptr) Uint8 {
	*c = Uint8(p)
	return *c
}

// Int16
func (c *Int16) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Int16) FromPtr(p uintptr) Int16 {
	*c = Int16(p)
	return *c
}

// Uint16
func (c *Uint16) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Uint16) FromPtr(p uintptr) Uint16 {
	*c = Uint16(p)
	return *c
}

// Int32
func (c *Int32) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Int32) FromPtr(p uintptr) Int32 {
	*c = Int32(p)
	return *c
}

// Uint32
func (c *Uint32) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Uint32) FromPtr(p uintptr) Uint32 {
	*c = Uint32(p)
	return *c
}

// Int64
func (c *Int64) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Int64) FromPtr(p uintptr) Int64 {
	*c = Int64(p)
	return *c
}

// Uint64
func (c *Uint64) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Uint64) FromPtr(p uintptr) Uint64 {
	*c = Uint64(p)
	return *c
}

// Float32
func (c *Float32) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Float32) FromPtr(p uintptr) Float32 {
	*c = Float32(p)
	return *c
}

// Float64
func (c *Float64) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Float64) FromPtr(p uintptr) Float64 {
	*c = Float64(p)
	return *c
}

// Bool
func (c *Bool) ToPtr() uintptr {
	return BoolToPtr(bool(*c))
}
func (c *Bool) FromPtr(p uintptr) Bool {
	*c = Bool(PtrToBool(p))
	return *c
}

// Uintptr
func (c *Uintptr) ToPtr() uintptr {
	return uintptr(*c)
}
func (c *Uintptr) FromPtr(p uintptr) Uintptr {
	*c = Uintptr(p)
	return *c
}

// String
func (c *String) ToPtr() uintptr {
	return StringToPtr(string(*c))
}
func (c *String) FromPtr(p uintptr) String {
	*c = String(PtrToString(p))
	return *c
}

// WString
func (c *WString) ToPtr() uintptr {
	return StringToWCharPtr(string(*c))
}
func (c *WString) FromPtr(p uintptr) WString {
	*c = WString(PtrWCharToString(p))
	return *c
}

// AnyPtr
func (c *AnyPtr) ToPtr() uintptr {
	return uintptr(unsafe.Pointer(c))
}
func (c *AnyPtr) FromPtr(p uintptr) AnyPtr {
	*c = *(*AnyPtr)(unsafe.Pointer(p))
	return *c
}
