package parser

import (
	_ "embed"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

//go:embed wke.include.h
var wkeInclude string

func removeBOM(content []byte) []byte {
	if len(content) >= 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		return content[3:]
	}
	return content
}

func Simplify(input string) ([]byte, error) {
	// 精简include
	byts, err := os.ReadFile(input)
	if err != nil {
		return nil, err
	}

	byts = removeBOM(byts)

	re := regexp.MustCompile(`\#include\s+["<](.*)[">]`)
	replaced := re.ReplaceAllString(string(byts), "")
	replaced = wkeInclude + "\n" + replaced

	// 创建临时目录
	var dir string
	dir, err = os.MkdirTemp("", "output")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	// 写入临时文件
	tempFile := path.Join(dir, "input.h")
	os.WriteFile(tempFile, []byte(replaced), 0644)

	// 预处理头文件
	output := path.Join(dir, "parsed.h")

	err = exec.Command("gcc", "-E", "-P", tempFile, "-o", output).Run()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(output)
	if err != nil {
		return nil, err
	}

	// 格式化
	str := string(content)
	str = strings.ReplaceAll(str, "_Bool", "Bool")
	str = strings.ReplaceAll(str, "; __declspec(dllexport)", ";\n __declspec(dllexport)")

	str = reWarpLines(str)

	str = fixAnonymousStruct(str)
	str = separeteStruct(str)
	str = separeteEnum(str)

	return []byte(str), nil
}

// 将匿名结构体转换为具名结构体
func fixAnonymousStruct(text string) string {
	re := regexp.MustCompile(`typedef\s+struct\s*{([^{}]*)}\s*(\w+);`)
	return re.ReplaceAllString(text, "struct _$2 { $1 }; \ntypedef struct _$2 $2;")
}

// 将结构体的定义和声明分开
func separeteStruct(text string) string {
	re := regexp.MustCompile(`typedef\s+struct\s+(\w+)\s*{([^{}]*)}\s*(\w+);`)
	return re.ReplaceAllString(text, "struct $1 { $2 }; \ntypedef struct $1 $3;")
}

// 将枚举的定义和声明分开
func separeteEnum(text string) string {
	re := regexp.MustCompile(`typedef\s+enum\s+(\w+)\s*{([^{}]*)}\s*(\w+);`)
	return re.ReplaceAllString(text, "enum $1 { $2 }; \ntypedef enum $1 $3;")
}

func reWarpLines(text string) string {
	// 移除所有换行
	// text = strings.ReplaceAll(text, "\r\n", " ")
	// text = strings.ReplaceAll(text, "\n", " ")

	// 根据分号分割换行
	re := regexp.MustCompile(`;\s*([^\s][^;\r\n(\/\/)]+)`) // 注释不换行
	text = re.ReplaceAllString(text, ";\n$1")

	// 大括号换行
	text = strings.ReplaceAll(text, "{", "{\n")
	text = strings.ReplaceAll(text, "}", "\n}")

	return text
}
