package parser

import (
	"regexp"
	"strings"
)

func Prettify(txt *string) {
	removeCPPBlocks(txt)
	replaceWKE_CALL_TYPE(txt)
}

// 移除C++预处理指令
func removeCPPBlocks(txt *string) {
	re := regexp.MustCompile(`\#if\s*defined\(__cplusplus\)[\s\S]*?\#endif`)
	*txt = re.ReplaceAllString(*txt, "")

	// 注意：这个正则表达式假设!defined#else之间没有嵌套的其他预处理指令
	re = regexp.MustCompile(`\#if\s*\!defined\(__cplusplus\)([\s\S]*?)\#else([\s\S]*?)\#endif`)

	// 使用$1来引用第一个捕获组，即!defined和#else之间的内容
	*txt = re.ReplaceAllString(*txt, `$1`)
}

// 替换 WKE_CALL_TYPE
func replaceWKE_CALL_TYPE(txt *string) {
	*txt = strings.Replace(*txt, "#define WKE_CALL_TYPE __cdecl", "", 1)
	re := regexp.MustCompile(`WKE_CALL_TYPE\s*`)
	*txt = re.ReplaceAllString(*txt, "__cdecl ")
}

var defs map[string]interface{} = map[string]interface{}{
	"WKE_CALL_TYPE": "__cdecl",
}

type IfBlock struct {
	Condition string
	Lines     []string
	ElseLines []string
}

// func pickIFBlock(txt *string) string {
// 	lines := strings.Split(*txt, "\n")

// 	result := ""

// 	for i, line := range lines {

// 		if strings.HasPrefix(line, "#if") {
// 			re := regexp.MustCompile(`\#if\s*(\!?)defined\((\w+)\)`)
// 			match := re.FindStringSubmatch(line)
// 			isDefined := match[1] != "!"
// 			constant := match[2]

// 			_, hasDef := defs[constant]
// 			isValid := isDefined == hasDef

// 		} else if result != "" && strings.Contains(line, "#endif") {
// 			break
// 		}
// 	}
// }
