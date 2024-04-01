package blink

import (
	"os"
	"path/filepath"
)

type Config struct {
	// 临时文件夹，用于释放 DLL 以及其他临时文件
	tempPath string
	// dll文件路径，非绝对路径将在临时文件夹内创建
	dll string
	// 设置storage本地文件目录
	storagePath string
	// 设置cookie文件名
	cookieFile string
}

func NewConfig(setups ...func(*Config)) *Config {

	conf := &Config{
		tempPath:    os.TempDir(),
		dll:         "blink_" + AppID + ".dll",
		storagePath: "LocalStorage_" + AppID,
		cookieFile:  "cookie_" + AppID + ".dat",
	}

	for _, setup := range setups {
		setup(conf)
	}

	return conf
}

func WithConfigTempPath(path string) func(*Config) {
	return func(conf *Config) {
		conf.tempPath = path
	}
}

func WithConfigDll(dll string) func(*Config) {
	return func(conf *Config) {
		conf.dll = dll
	}
}

func WithConfigStoragePath(path string) func(*Config) {
	return func(conf *Config) {
		conf.storagePath = path
	}
}

func WithConfigCookieFile(path string) func(*Config) {
	return func(conf *Config) {
		conf.cookieFile = path
	}
}

func (conf *Config) GetDllFilePath() string {

	if filepath.IsAbs(conf.dll) {
		return conf.dll
	}

	return filepath.Join(conf.tempPath, conf.dll)
}

func (conf *Config) GetStoragePath() string {

	if filepath.IsAbs(conf.storagePath) {
		return conf.storagePath
	}

	return filepath.Join(conf.tempPath, conf.storagePath)
}

func (conf *Config) GetCookieFilePath() string {

	if filepath.IsAbs(conf.cookieFile) {
		return conf.cookieFile
	}

	return filepath.Join(conf.tempPath, conf.cookieFile)
}
