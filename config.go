package blink

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/epkgs/blink/internal/log"
)

type Config struct {
	// 临时文件夹，用于释放 DLL 以及其他临时文件
	tempPath string
	// dll文件路径，非绝对路径将在临时文件夹内创建
	dllFile string
	// 设置storage本地文件目录
	storagePath string
	// 设置cookie文件名
	cookieFile string
}

func NewConfig(setups ...func(*Config)) (*Config, error) {

	tempPath := filepath.Join(os.TempDir(), "mini-blink")

	conf := &Config{
		tempPath:    tempPath,
		dllFile:     "blink.dll",
		storagePath: "LocalStorage",
		cookieFile:  "cookie.dat",
	}

	for _, setup := range setups {
		setup(conf)
	}

	log.Debug("临时文件夹：%s", conf.tempPath)
	if err := os.MkdirAll(conf.tempPath, 0644); err != nil {
		return nil, fmt.Errorf("临时文件夹(%s)不存在，且创建不成功，请确认文件夹权限。", conf.tempPath)
	}

	return conf, nil
}

func WithTempPath(path string) func(*Config) {
	return func(conf *Config) {
		conf.tempPath = path
	}
}

func WithDllFile(dllFile string) func(*Config) {
	return func(conf *Config) {
		conf.dllFile = dllFile
	}
}

func WithStoragePath(path string) func(*Config) {
	return func(conf *Config) {
		conf.storagePath = path
	}
}

func WithCookieFile(path string) func(*Config) {
	return func(conf *Config) {
		conf.cookieFile = path
	}
}

func (conf *Config) GetDllFile() string {
	return conf.dllFile
}

func (conf *Config) GetTempPath() string {
	return conf.tempPath
}

func (conf *Config) GetDllFileABS() string {

	if filepath.IsAbs(conf.dllFile) {
		return conf.dllFile
	}

	return filepath.Join(conf.tempPath, conf.dllFile)
}

func (conf *Config) GetStoragePath() string {

	if filepath.IsAbs(conf.storagePath) {
		return conf.storagePath
	}

	return filepath.Join(conf.tempPath, conf.storagePath)
}

func (conf *Config) GetCookieFileABS() string {

	if filepath.IsAbs(conf.cookieFile) {
		return conf.cookieFile
	}

	return filepath.Join(conf.tempPath, conf.cookieFile)
}
