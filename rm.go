package goshell

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// 删除文件或者目录
func Rm(filePath string) error {
	return os.Remove(filePath)
}

// 移动
func Move(oldpath, newpath string) error {
	fmt.Println("move", oldpath, newpath)
	return os.Rename(oldpath, newpath)
}

// 移动
func MustMove(oldpath, newpath string) {
	err := Move(oldpath, newpath)
	assert("rename "+oldpath+" -> "+newpath, err)
}

// 改名字
func Rename(oldpath, newpath string) error {
	fmt.Println("rename", oldpath, newpath)
	return os.Rename(oldpath, newpath)
}

// 改名字
func MustRename(oldpath, newpath string) {
	err := Rename(oldpath, newpath)
	assert("rename "+oldpath+" -> "+newpath, err)
}

// 清空文件夹
func ClearDir(dirPath string) error {
	items, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, d := range items {
		err := os.RemoveAll(path.Join([]string{"tmp", d.Name()}...))
		if err != nil {
			return err
		}
	}
	return nil
}
