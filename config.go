package blink

import (
	"os"
	"path/filepath"
)

type Config struct {
	// 32位dll文件名
	dll32 string
	// 64位dll文件名
	dll64 string
	// (运行文件夹 / 临时文件夹) 注意：需要以 / 结尾
	runtimePath string
	// 设置cookie的本地文件目录。默认是当前目录。cookies存在当前目录的“cookie.dat”里
	storagePath string
	// 设置cookie的全路径+文件名，如“c:\mb\cookie.dat”
	cookieFile string
}

func NewConfig() *Config {

	runtimePath, err := os.Getwd()
	if err != nil {
		runtimePath = "."
	}

	err = os.MkdirAll(runtimePath, 0644)
	if err != nil {
		pwd, err := os.Getwd()
		if err != nil {
			panic("运行目录无权限！")
		}
		runtimePath = pwd
	}

	return &Config{
		dll32:       "blink_x32.dll",
		dll64:       "blink_x64.dll",
		storagePath: "LocalStorage",
		cookieFile:  "cookie.dat",
		runtimePath: runtimePath,
	}
}

func (conf *Config) GetDllFilePath() string {

	dll := conf.dll32
	if env.isSYS64 {
		dll = conf.dll64
	}

	if filepath.IsAbs(dll) {
		return dll
	}

	return filepath.Join(conf.runtimePath, dll)
}

func (conf *Config) GetStoragePath() string {

	if filepath.IsAbs(conf.storagePath) {
		return conf.storagePath
	}

	return filepath.Join(conf.runtimePath, conf.storagePath)
}

func (conf *Config) GetCookieFilePath() string {

	if filepath.IsAbs(conf.cookieFile) {
		return conf.cookieFile
	}

	return filepath.Join(conf.runtimePath, conf.cookieFile)
}
