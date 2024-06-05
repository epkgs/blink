package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type CEnumField struct {
	Name  string
	Value string
}

type CEnum struct {
	Name   string
	Fields []CEnumField
}

func (c *CEnumField) GoName() string {
	return UpperFirst(c.Name)
}

func (c *CEnum) GoName() string {
	return UpperFirst(c.Name)
}

func (c *CEnum) Placeholder() string {
	return fmt.Sprintf("[enum:%s]", c.Name)
}

func ParseEnums(str *string) []CEnum {
	var result []CEnum

	re := regexp.MustCompile(`typedef\s+enum\s+_(\w+)\s*\{([\s\S]*?)\}\s*(\w+)\s*[;\r\n]`)

	matches := re.FindAllStringSubmatch(*str, -1)

	for _, match := range matches {
		var item CEnum
		if match[1] != match[3] {
			continue
		}

		item.Name = match[1]
		*str = strings.Replace(*str, match[0], item.Placeholder(), 1)

		fieldRE := regexp.MustCompile(`\s*(\w+)\s+[=\s]?(.*)\s*[,\r\n]`)
		fieldMatches := fieldRE.FindAllStringSubmatch(match[2], -1)
		fields := []CEnumField{}
		for _, fm := range fieldMatches {
			if fm[1] == "" {
				continue
			}

			fields = append(fields, CEnumField{Name: fm[1], Value: fm[2]})
		}
		item.Fields = fields

		result = append(result, item)
	}

	return result
}
