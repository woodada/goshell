package goshell

import (
	"io/ioutil"
	"os"
)

func Ls(dirPath string) (files []string, dirs []string, err error) {
	var arr []os.FileInfo
	arr, err = ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	for _, a := range arr {
		if a.IsDir() {
			dirs = append(dirs, a.Name())
		} else {
			files = append(files, a.Name())
		}
	}
	return
}
