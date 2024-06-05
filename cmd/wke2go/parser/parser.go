package parser

import (
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

//go:embed output.go.tpl
var dllTPL string

type TplData struct {
	PkgName     string
	FuncMapping map[string]string
}

func Parse(pkgName string, text string) string {

	result := "\n"

	lines := breakLines(text)
	maxIdx := len(lines) - 1
	i := 0
	for i <= maxIdx {

		if pickCallbackDef(&result, &lines, &i) {
			continue
		}

		if pickStruct(&result, &lines, &i) {
			continue
		}

		if pickEnum(&result, &lines, &i) {
			continue
		}

		if pickTypeDef(&result, &lines, &i) {
			continue
		}

		if pickFuncExport(&result, &lines, &i) {
			continue
		}

		i++
	}

	data := TplData{
		PkgName:     pkgName,
		FuncMapping: FuncMapping,
	}

	var buf bytes.Buffer

	templ := template.Must(template.New("output").Parse(dllTPL))

	templ.Execute(&buf, data)

	buf.Write([]byte("\n\n"))
	buf.Write([]byte(generateGoTypes()))
	buf.Write([]byte("\n\n"))
	buf.Write([]byte(result))

	return buf.String()
}

func pickTypeDef(result *string, lines *[]string, idx *int) bool {
	line := (*lines)[*idx]

	re := regexp.MustCompile(`^\s*typedef\s+((struct|enum))?\s*([\w\s\*]+)\s+\*?(\w+);`)
	match := re.FindStringSubmatch(line)

	if match == nil {
		return false
	}

	isStruct := match[2] == "struct"
	isEnum := match[2] == "enum"
	CType := match[3]
	CName := match[4]
	GName := UpperFirst(CName)

	prevType := GetGType(CType)

	if prevType == "" { // 未定义

		if isStruct {
			re = regexp.MustCompile(`\s*(\w+)\s*`)
			match = re.FindStringSubmatch(CType)
			ctype := match[1]
			*result += fmt.Sprintf("\ntype %s struct {}\n", ctype)
			AddTypeMapping(CName, GName)
			AddTypeDef(GName, ctype)
		} else if isEnum {
			AddTypeMapping(CName, GName)
		} else {
			// 使用预定义的 Any 类型
			AddTypeMapping(CName, "AnyPtr")
		}

	} else { // 已定义

		if isStruct {
			AddTypeMapping(CName, GName)
			AddTypeDef(GName, prevType)
		} else if isEnum {
			// 替换原来的ENUM类型定义
			*result = strings.Replace(*result, prevType, GName, -1)
			AddTypeMapping(CName, GName)
		} else {

			// if CType != CName {
			// 	AddTypeMapping(CName, GName)
			// }
		}

	}

	(*idx)++

	return true
}

var callbackDefTPL = `
// %s
type %s func(%s) (%s)
func(cb *%s) ToPtr() uintptr {
	return CallbackToPtr(*cb)
}
func(cb *%s) FromPtr(p uintptr) %s{
	panic("暂未实现从 uintptr 转换为 %s")
}

`

func pickCallbackDef(result *string, lines *[]string, idx *int) bool {
	line := (*lines)[*idx]
	re := regexp.MustCompile(`typedef\s+(\w+\*?)\((\w+)\s*\*\s*(\w+)\)\s*\(([^\)]*)(\)\s*;)?`)
	match := re.FindStringSubmatch(line)
	if match == nil {
		return false
	}

	if match[5] == "" {
		// 回调定义不完整
		for {
			(*idx)++
			line += (*lines)[*idx]
			match = re.FindStringSubmatch(line)
			if match[5] != "" {
				break
			}
		}
	}

	original := match[0]
	CReturn := match[1]
	// CExportType:=	match[2]
	CName := match[3]
	GName := UpperFirst(CName)
	CArgStr := match[4]
	GArgs := []string{} // []"GName GBaseType"
	if CArgStr != "" {
		re = regexp.MustCompile(`\s*(const)?\s*(struct|enum)?\s*(\w+)\s*(\*)?\s*([a-zA-Z0-9_]+)\s*,?`)
		match := re.FindAllStringSubmatch(CArgStr, -1)
		for _, m := range match {
			ctype := m[3] + m[4] // 可能带 * 号
			name := m[5]

			if isKeyWord(name) {
				name += "_"
			}

			isStruct := m[2] == "struct"
			// isEnum := m[2] == "enum"

			gtype := GetGType(ctype)

			if gtype == "" {
				panic("pickCallbackDef, 回调 " + CName + " 参数未知类型: " + ctype)
			}

			ref := ""
			if isStruct {
				ref = "*"
			}

			GArgs = append(GArgs, name+" "+ref+gtype)
		}
	}

	var GReturn string
	if CReturn == "void" {
		GReturn = "void uintptr"
	} else {
		gtype := GetGType(CReturn)
		if gtype == "" {
			panic("pickCallbackDef, 回调 " + CName + " 返回值未知类型: " + CReturn)
		}
		GReturn = gtype
	}

	*result += fmt.Sprintf(callbackDefTPL, original, GName, strings.Join(GArgs, ", "), GReturn, GName, GName, GName, GName)
	AddTypeMapping(CName, GName)

	(*idx)++
	return true
}

var funcVoidTPL = `
// %s
func %s(%s) (err error) {
	_, _, err = globalDLL._%s.Call(%s)
	return
}

`

var funcTPL = `
// %s
func %s(%s) (res %s, err error) {
	var r uintptr
	r, _, err = globalDLL._%s.Call(%s)

	if err != nil {
		return
	}

	res.FromPtr(r)
	return
}

`

func pickFuncExport(result *string, lines *[]string, idx *int) bool {
	line := (*lines)[*idx]

	re := regexp.MustCompile(`__declspec\(dllexport\)\s*(const)?\s*(struct|enum)?\s*([\w\s\*]+?)\s*__cdecl\s*(\w+)\((.*)\);`)
	match := re.FindStringSubmatch(line)
	if match == nil {
		return false
	}
	original := match[0]
	CReturn := match[3]
	CName := match[4]
	GName := UpperFirst(CName)
	CArgStr := match[5]
	GArgs := []string{}    // []"GName GBaseType"
	GArgPtrs := []string{} // []varible.ToPtr()

	if CArgStr != "" {
		re = regexp.MustCompile(`\s*(const)?\s*(struct|enum)?\s*(\w+[\w\s]*?)\s*([\s\*]+?)\s*(\w+)(,\s*)?`)
		match := re.FindAllStringSubmatch(CArgStr, -1)
		for _, m := range match {
			ctype := m[3] + m[4] // 可能带 * 号
			name := m[5]

			if isKeyWord(name) {
				name += "_"
			}

			isStruct := m[2] == "struct"
			// isEnum := m[2] == "enum"

			gtype := GetGType(ctype)

			if gtype == "" {
				panic("pickFuncExport, 函数 " + CName + " 参数未知类型: " + ctype)
			}

			ref := ""
			if isStruct {
				ref = "*"
			}

			GArgs = append(GArgs, name+" "+ref+gtype)
			GArgPtrs = append(GArgPtrs, name+".ToPtr()")
		}
	}

	var GReturn string
	if CReturn == "void" {
		*result += fmt.Sprintf(funcVoidTPL, original, GName, strings.Join(GArgs, ", "), CName, strings.Join(GArgPtrs, ", "))

	} else {
		gtype := GetGType(CReturn)
		if gtype == "" {
			panic("pickCallbackDef, 函数 " + CName + " 返回值未知类型: " + CReturn)
		}
		GReturn = gtype
		*result += fmt.Sprintf(funcTPL, original, GName, strings.Join(GArgs, ", "), GReturn, CName, strings.Join(GArgPtrs, ", "))
	}

	AddFuncMapping(CName, GName)

	(*idx)++
	return true
}

var structTPL = `
type %s struct {
	%s
}
func(s *%s) ToPtr() uintptr {
	return uintptr(unsafe.Pointer(s))
}
func(s *%s) FromPtr(p uintptr) %s{
	*s = *(*%s)(unsafe.Pointer(s))
	return *s
}

`

func pickStruct(result *string, lines *[]string, idx *int) bool {
	line := (*lines)[*idx]
	re := regexp.MustCompile(`^\s*struct\s+(\w+)\s*([{|;])`)
	match := re.FindStringSubmatch(line)
	if match == nil {
		return false
	}

	defer func() {
		(*idx)++
	}()

	CType := match[1]
	GType := CType

	AddTypeMapping(CType, GType) // 会自动覆盖原来的定义
	AddTypeDef(GType, GType)     // struct 本身就是一个类型，当等于它自身时，不再重复生成 type 定义

	if match[2] == ";" {
		AddTypeDef(GType, "AnyPtr")
		return true
	}

	propsStr := ""

	for {
		(*idx)++
		line = (*lines)[*idx]

		re := regexp.MustCompile(`\s*(\})\s*;`)
		match := re.FindStringSubmatch(line)
		if match != nil {
			break
		}

		// 正则表达式
		// \s+ 匹配一个或多个空白字符
		// (struct|enum)? 匹配 "struct" 或 "enum"，但它是可选的（因为有问号）
		// \s* 匹配零个或多个空白字符
		// (\w+\*?) 匹配一个或多个单词字符，后面可能跟着一个星号（因为星号是可选的，由问号表示）
		// \s+ 匹配一个或多个空白字符
		// (\w+) 匹配一个或多个单词字符
		// ; 匹配分号
		re = regexp.MustCompile(`\s*(struct|enum)?\s*(\w+\*?)\s+(\w+)(\[\d*?\])?;`)
		match = re.FindStringSubmatch(line)

		if match == nil {
			continue
		}

		isStruct := match[1] == "struct"
		// isEnum := match[1] == "enum"
		propCType := match[2]
		propCName := UpperFirst(match[3])
		quote := match[4]

		propGType := GetGType(propCType)

		if propGType == "" {
			panic("typeDef not found: " + propCType)
		}

		ref := ""
		if isStruct {
			ref = "*"
		}

		propsStr += "  " + propCName + " " + ref + quote + propGType + "\n"
	}

	*result += fmt.Sprintf(structTPL, GType, propsStr, GType, GType, GType, GType)

	return true
}

func pickEnum(result *string, lines *[]string, idx *int) bool {
	line := (*lines)[*idx]
	re := regexp.MustCompile(`enum\s+(\w+)\s*{`)
	match := re.FindStringSubmatch(line)
	if match == nil {
		return false
	}

	CType := match[1]
	GType := CType

	AddTypeMapping(CType, GType) // 会自动覆盖原来的定义

	*result += "\ntype " + GType + " Int\n"
	*result += "const (\n"

	firstLoop := true
	for {
		(*idx)++
		line = (*lines)[*idx]
		re := regexp.MustCompile(`\s*\}\s*;`)
		match := re.FindStringSubmatch(line)
		if match != nil {
			break
		}

		re = regexp.MustCompile(`\s*(\w+)(\s*\=?\s*)(.*),`)
		match = re.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		propCName := match[1]
		propGName := UpperFirst(propCName)
		// assert := match[2]
		val := match[3]
		if val == "" {
			if firstLoop {
				firstLoop = false
				*result += "  " + propGName + " " + GType + " = iota\n"
			} else {
				*result += "  " + propGName + "\n"
			}
		} else {
			if firstLoop {
				firstLoop = false
			}

			*result += fmt.Sprintf("  %s %s = %s\n", propGName, GType, val)
		}

	}

	*result += ")\n"

	(*idx)++

	return true
}

func generateGoTypes() string {

	definition := "\n"

	for gtype, baseType := range Typedefs {
		if gtype != baseType && baseType != "" {
			definition += "type " + gtype + " " + baseType + "\n"
		}
	}

	return definition
}

func breakLines(text string) []string {

	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")

	return lines

}

var GoKeyWords = []string{
	// 包管理
	"import", "package",
	// 程序实体声明与定义
	"chan", "const", "func", "interface", "map", "struct", "type", "var",
	// 程序流程控制
	"break", "case", "continue", "default", "defer", "else", "fallthrough",
	"for", "go", "goto", "if", "range", "return", "select", "switch",
}

func isKeyWord(name string) bool {
	for _, word := range GoKeyWords {
		if name == word {
			return true
		}
	}

	return false
}
