package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// 检查文件是否存在，并返回一个未使用的文件路径
func GetUnusedPath(originalPath string) string {

	// 检查文件是否存在
	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		// 文件不存在，返回新路径
		return originalPath
	}

	base := filepath.Base(originalPath)
	dir := filepath.Dir(originalPath)
	ext := filepath.Ext(base)
	baseWithoutExt := base[:len(base)-len(ext)]

	index := 1
	for {
		// 构造新的文件名
		newBase := fmt.Sprintf("%s(%d)%s", baseWithoutExt, index, ext)
		newPath := filepath.Join(dir, newBase)

		// 检查文件是否存在
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			// 文件不存在，返回新路径
			return newPath
		}

		// 文件存在，增加索引并重试
		index++
	}
}
