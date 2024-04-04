package blink

import (
	"errors"
	"io"
	"os"

	"github.com/epkgs/mini-blink/internal/dll"
	"golang.org/x/sys/windows"
)

func loadDLL(conf *Config) (*windows.DLL, error) {

	// 尝试在默认目录里加载 DLL
	if loaded, err := windows.LoadDLL(DLL_FILE); err == nil {
		return loaded, nil
	}

	fullPath := conf.GetDllFilePath()

	// 放入闭包，使其可以被释放
	file, err := dll.FS.Open(DLL_FILE)
	if err != nil {
		return nil, errors.New("无法从默认路径或内嵌资源里找到 blink.dll，err: " + err.Error())
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("读取内联DLL出错，err: " + err.Error())
	}

	newFile, err := os.Create(fullPath)
	if err != nil {
		return nil, errors.New("无法创建dll文件，err: " + err.Error())
	}
	defer newFile.Close()

	n, err := newFile.Write(data)
	if err != nil {
		return nil, errors.New("写入dll文件失败，err: " + err.Error())
	}
	if n != len(data) {
		return nil, errors.New("写入校验失败")
	}

	return windows.MustLoadDLL(fullPath), nil
}
