package parser

import "strings"

// UpperFirst 函数将一个字符串的首字母转换为大写，并返回转换后的字符串。
// 如果字符串为空或者首字符不是字母，则直接返回原字符串。
// 参数s：待转换的字符串。
// 返回值：转换后的字符串。
func UpperFirst(s string) string {
	if s == "" {
		return ""
	}

	return strings.ToUpper(s[0:1]) + s[1:]
}
