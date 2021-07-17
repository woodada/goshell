package goshell

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

var _config Config

func init() {
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "loadConfig", err.Error())
		return
	}
	for k, v := range config.Which {
		SetWhich(k, v)
	}
	_config = config
}

type Config struct {
	Version     string            `yml:"version"`
	Which       map[string]string `yml:"which"`
	DisableMuti bool              `yml:"disableMuti"`
}

// 文件存在才打开
func tryOpenFile(filePath string) (*os.File, error) {
	fileInfo := GetFileInfo(filePath)
	if fileInfo != nil && !fileInfo.IsDir() {
		return os.OpenFile(filePath, os.O_RDONLY, 0644)
	}
	return nil, nil
}

func IsStrInArr(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func loadConfig() (Config, error) {
	const cfgName = "autobuild.yml"
	var cfgFile *os.File
	var cfgPath string
	var tryPaths []string

	if cfgFile == nil {
		curDir, err := os.Getwd()
		if err == nil {
			cfgPath = filepath.Join(curDir, cfgName)
			if !IsStrInArr(cfgPath, tryPaths) {
				tryPaths = append(tryPaths, cfgPath)
				f, err := tryOpenFile(cfgPath)
				if err != nil {
					if err != nil {
						fmt.Fprintln(os.Stderr, "os.OpenFile", cfgPath, "err:", err)
						return Config{}, err
					}
				}
				cfgFile = f
				defer cfgFile.Close()
			}
		} else {
			fmt.Fprintln(os.Stderr, "os.Getwd err:", err)
		}
	}

	if cfgFile == nil {
		exePath, err := os.Executable()
		if err == nil {
			cfgPath = filepath.Join(filepath.Dir(exePath), cfgName)
			if !IsStrInArr(cfgPath, tryPaths) {
				tryPaths = append(tryPaths, cfgPath)
				f, err := tryOpenFile(cfgPath)
				if err != nil {
					fmt.Fprintln(os.Stderr, "os.OpenFile", cfgPath, "err:", err)
					return Config{}, err
				}
				cfgFile = f
				defer cfgFile.Close()
			}
		} else {
			fmt.Fprintln(os.Stderr, "os.Executable err:", err)
		}
	}

	if cfgFile != nil {
		var config Config
		err := yaml.NewDecoder(cfgFile).Decode(&config)
		if err == nil {
			fmt.Println("[加载配置] 加载", cfgPath, "成功")
		} else {
			fmt.Println("[加载配置] 加载", cfgPath, "失败", err.Error())
		}
		return config, err
	}

	fmt.Println("[加载配置] 未找到配置文件", strings.Join(tryPaths, " "))
	return Config{}, nil
}
