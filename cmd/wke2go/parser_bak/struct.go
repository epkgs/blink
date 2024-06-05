package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type CStructField struct {
	Name string
	Type string
}

type CStruct struct {
	Name   string
	Fields []CStructField
}

func (c *CStructField) GoName() string {
	return UpperFirst(c.Name)
}

func (c *CStruct) GoName() string {
	return UpperFirst(c.Name)
}

func (c *CStruct) Placeholder() string {
	return fmt.Sprintf("[struct:%s]", c.Name)
}

func ParseStructs(str *string) []CStruct {
	var result []CStruct

	// 正则表达式匹配结构体的名称、属性名和类型
	re := regexp.MustCompile(`typedef struct\s+_(\w+)\s*\{([\s\S]*?)\}\s*(\w+)\s*[;\r\n]`)
	matches := re.FindAllStringSubmatch(*str, -1)

	for _, match := range matches {
		var item CStruct
		if match[1] != match[3] {
			continue
		}

		item.Name = match[1]

		*str = strings.Replace(*str, match[0], item.Placeholder(), 1)

		fieldRE := regexp.MustCompile(`\s*(\w+)\s+(\w+)\s*[;\r\n]`)
		fieldMatches := fieldRE.FindAllStringSubmatch(match[2], -1)
		fields := []CStructField{}
		for _, fm := range fieldMatches {
			if fm[1] == "" || fm[2] == "" {
				continue
			}
			fields = append(fields, CStructField{Name: fm[2], Type: fm[1]})
		}
		item.Fields = fields

		result = append(result, item)
	}

	return result
}
