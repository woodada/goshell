package goshell

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 判断文件或目录是否存在
func GetFileInfo(src string) os.FileInfo {
	if fileInfo, e := os.Stat(src); e != nil {
		if os.IsNotExist(e) {
			return nil
		}
		return nil
	} else {
		return fileInfo
	}
}

// 拷贝文件
func CopyFile(src, dst string) error {
	if len(src) == 0 || len(dst) == 0 {
		return fmt.Errorf("EmptyPath")
	}
	srcFile, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFileInfo := GetFileInfo(dst)
	if dstFileInfo == nil {
		if e := os.MkdirAll(filepath.Dir(dst), os.ModePerm); e != nil {
			return e
		}
	}
	//这里要把O_TRUNC 加上，否则会出现新旧文件内容出现重叠现象
	dstFile, e := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if e != nil {
		return e
	}
	defer dstFile.Close()
	//fileInfo, e := srcFile.Stat()
	//fileInfo.Size() > 1024
	//byteBuffer := make([]byte, 10)
	if _, e := io.Copy(dstFile, srcFile); e != nil {
		return e
	}
	return nil
}

// 拷贝目录
func CopyDir(src, dst string) error {
	var err error
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}

	srcFileInfo := GetFileInfo(src)
	if srcFileInfo == nil || !srcFileInfo.IsDir() {
		return fmt.Errorf("%sNotExists", src)
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, walkerr error) error {
		if walkerr != nil {
			return walkerr
		}
		dstPath := filepath.Join(dst, strings.TrimLeft(path, src))
		// 是文件
		if !info.IsDir() {
			return CopyFile(path, dstPath)
		}
		// 目录
		_, e := os.Stat(dstPath)
		if e != nil && os.IsNotExist(e) {
			return os.MkdirAll(dstPath, os.ModePerm)
		}
		return e
	})

	return err
}
