package parser

import (
	"regexp"
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

// GType => Base Type
var Typedefs = map[string]string{
	"Size_t":                "Uintptr",
	"WCHAR":                 "WString",
	"DWORD":                 "UInt",
	"WORD":                  "UShort",
	"HANDLE":                "Uintptr",
	"LPWSTR":                "WString",
	"PWSTR":                 "WString",
	"BYTE":                  "Byte",
	"LPBYTE":                "Byte",
	"LONG":                  "Long",
	"COLORREF":              "Uint32",
	"BOOL":                  "Bool",
	"UINT":                  "Uint32",
	"WPARAM":                "Uint32",
	"LPARAM":                "Int32",
	"LRESULT":               "Int32",
	"JsExecState":           "Uintptr",
	"JsValue":               "Uintptr",
	"WkeWebView":            "Uintptr",
	"WkeString":             "Uintptr",
	"WkeMediaPlayer":        "Uintptr",
	"WkeMediaPlayerClient":  "Uintptr",
	"BlinkWebURLRequestPtr": "Uintptr",
	"WkeWebFrameHandle":     "Uintptr",
	"WkeNetJob":             "Uintptr",
}

// CType => GType
var TypeMapping = map[string]string{
	// 基本类型
	"char":               "Char",      // byte
	"singed char":        "SChar",     // int8
	"unsigned char":      "UChar",     // uint8
	"short":              "Short",     // int16
	"unsigned short":     "UShort",    // uint16
	"int":                "Int",       // int32
	"unsigned int":       "UInt",      // uint32
	"long":               "Long",      // int32
	"unsigned long":      "ULong",     // uint32
	"long long":          "LongLong",  // int64
	"unsigned long long": "ULongLong", // uint64
	"float":              "Float",     // float32
	"double":             "Double",    // float64
	"Bool":               "Bool",      // bool
	"void*":              "AnyPtr",    // uintptr

	// 字符串指针类型
	"utf8":    "String",
	"wchar_t": "WString",

	// 自定义类型
	"unsigned":              "UInt",
	"size_t":                "Size_t",
	"int64":                 "Int64",
	"WCHAR":                 "WCHAR",
	"DWORD":                 "DWORD",
	"WORD":                  "WORD",
	"HANDLE":                "HANDLE",
	"LPWSTR":                "LPWSTR",
	"PWSTR":                 "PWSTR",
	"BYTE":                  "BYTE",
	"LPBYTE":                "LPBYTE",
	"LONG":                  "LONG",
	"COLORREF":              "COLORREF",
	"BOOL":                  "BOOL",
	"UINT":                  "UINT",
	"WPARAM":                "WPARAM",
	"LPARAM":                "LPARAM",
	"LRESULT":               "LRESULT",
	"jsExecState":           "JsExecState",
	"jsValue":               "JsValue",
	"wkeWebView":            "WkeWebView",
	"wkeString":             "WkeString",
	"wkeMediaPlayer":        "WkeMediaPlayer",
	"wkeMediaPlayerClient":  "WkeMediaPlayerClient",
	"blinkWebURLRequestPtr": "BlinkWebURLRequestPtr",
	"wkeWebFrameHandle":     "WkeWebFrameHandle",
	"wkeNetJob":             "WkeNetJob",
}

var FuncMapping = map[string]string{}

// 会自动覆盖原来的定义
func AddTypeDef(gtype string, basetype string) {
	Typedefs[gtype] = basetype
}

// 会自动覆盖原来的定义
func AddTypeMapping(ctype string, gtype string) {
	TypeMapping[ctype] = gtype
}

func AddFuncMapping(cname string, gname string) {
	FuncMapping[cname] = gname
}

func GetGType(ctype string) string {

	re := regexp.MustCompile(`^\s*(.*?)\s*(\*?)\s*$`)
	match := re.FindStringSubmatch(ctype)
	if match == nil {
		return ""
	}
	if match[1] == "void" && match[2] == "*" {
		return "AnyPtr"
	}

	if gtype, ok := TypeMapping[match[1]]; ok {
		return gtype
	}

	return ""
}

func GetCType(gtype string) string {

	for c, g := range TypeMapping {
		if g == gtype {
			return c
		}
	}

	return ""
}
