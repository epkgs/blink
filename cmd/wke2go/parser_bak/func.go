package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type CFuncArg struct {
	Name string // 参数名
	Type string // 参数类型
}

type CFunc struct {
	Name     string     // 函数名
	Args     []CFuncArg // [][参数名，参数类型]
	Result   string     // 返回值类型
	Comments []string   // 注释
}

func (f *CFunc) GoName() string {
	return UpperFirst(f.Name)
}

func (c *CFunc) Placeholder() string {
	return fmt.Sprintf("[func:%s]", c.Name)
}

func ParseFuncs(txt *string) []CFunc {
	original := pickFuncBlock(*txt)
	replaced := formatFuncBlock(original)
	funcs := pickFuncPices(&replaced)

	*txt = strings.Replace(*txt, original, replaced, 1)
	re := regexp.MustCompile(`#define WKE_FOR_EACH_DEFINE_FUNCTION\([^)]*\) *\\`)
	*txt = re.ReplaceAllString(*txt, "")

	return funcs
}

func pickFuncBlock(txt string) string {
	re := regexp.MustCompile(`(?s)#define WKE_FOR_EACH_DEFINE_FUNCTION\([^)]*\) *\\(.*)#if ENABLE_WKE == 1`)

	match := re.FindStringSubmatch(txt)

	return match[1]
}

func formatFuncBlock(funcBlock string) string {
	// 移除无用的换行
	re := regexp.MustCompile(`\) *\\\r?\n`)
	replaced := re.ReplaceAllString(funcBlock, ")\n")

	re = regexp.MustCompile(` *\\\r?\n`)
	replaced = re.ReplaceAllString(replaced, "")

	return replaced
}

func pickFuncPices(str *string) []CFunc {
	re := regexp.MustCompile(`ITERATOR\d{1,2}\(([^,]*), ([^,]*), ([^"]*) *((\"([^"]*)\" *)+)\)`) // 返回值， 参数名， 参数， 注释

	matchs := re.FindAllStringSubmatch(*str, -1)

	funcs := make([]CFunc, len(matchs))
	for i, match := range matchs {
		funcs[i].Result = match[1]
		funcs[i].Name = match[2]

		// 替换占位符
		*str = strings.Replace(*str, match[0], funcs[i].Placeholder(), 1)

		args := match[3]
		args = strings.TrimRight(args, ",") + ","
		argsRE := regexp.MustCompile(`([^ ]*) ([^ ,]*),`)
		argMatchs := argsRE.FindAllStringSubmatch(args, -1)
		argsArr := []CFuncArg{}
		for _, argMatch := range argMatchs {
			if argMatch[1] == "" && argMatch[2] == "" {
				continue
			}

			argsArr = append(argsArr, CFuncArg{argMatch[2], argMatch[1]})
		}

		funcs[i].Args = argsArr

		comments := match[4]
		commentsRE := regexp.MustCompile(`\"([^"]+)\"`)
		commentMatchs := commentsRE.FindAllStringSubmatch(comments, -1)
		commentsArr := []string{}
		for _, commentMatch := range commentMatchs {
			if commentMatch[1] == "" {
				continue
			}
			commentsArr = append(commentsArr, commentMatch[1])
		}
		funcs[i].Comments = commentsArr
	}

	return funcs
}
