package miniblink

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/epkgs/blink/internal/log"
	"golang.org/x/sys/windows"
)

func LoadDLL(dllFile, tempPath string) (*windows.DLL, error) {
	// 建议将dll做个hash，加载时批判默认目录的dll是否为完整的dll，不完整的话从小释放dll
	// 尝试直接从默认目录里加载 DLL
	if loaded, err := windows.LoadDLL(dllFile); err == nil {
		log.Debug("直接加载DLL: %s", dllFile)
		return loaded, nil
	}

	dir := filepath.Join(tempPath, fmt.Sprintf("miniblink_%s_%s", VERSION, ARCH))
	releaseFile := filepath.Join(dir, dllFile)

	// 尝试直接加载释放后的 DLL
	if loaded, err := windows.LoadDLL(releaseFile); err == nil {
		log.Debug("直接加载DLL: %s", releaseFile)
		return loaded, nil
	}

	// 没有找到释放后的 DLL，则尝试从内嵌资源里释放

	// 临时文件夹不存在，则创建
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, errors.New("无法创建临时文件夹，err: " + err.Error())
	}

	// 释放内嵌资源里的 DLL
	if err := releaseEmbedDLL(releaseFile); err != nil {
		return nil, err
	}

	log.Debug("从内嵌资源里释放并加载 %s", releaseFile)
	return windows.MustLoadDLL(releaseFile), nil
}

func releaseEmbedDLL(releaseFile string) error {

	// 尝试从内嵌资源里打开 DLL 文件
	file, err := res.Open(fmt.Sprintf("release/%s/miniblink_%s_%s.dll", ARCH, VERSION, ARCH))
	if err != nil {
		return errors.New("无法从默认路径或内嵌资源里找到 blink.dll，err: " + err.Error())
	}
	defer file.Close()

	// 读取内嵌资源 DLL 文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		return errors.New("读取内联DLL出错，err: " + err.Error())
	}

	// 创建dll文件
	newFile, err := os.Create(releaseFile)
	if err != nil {
		return errors.New("无法创建dll文件，err: " + err.Error())
	}
	defer newFile.Close()

	n, err := newFile.Write(data)
	if err != nil {
		return errors.New("写入dll文件失败，err: " + err.Error())
	}
	if n != len(data) {
		return errors.New("写入校验失败")
	}

	return nil
}
