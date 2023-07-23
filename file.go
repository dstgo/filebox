package filebox

import (
	"errors"
	"os"
	"path"
)

// CreateFile
// param file string
// return *os.File
// return error
// 创建一个指定名称的文件，并且会检查文件的父目录是否存在
func CreateFile(file string) (*os.File, error) {
	dir := path.Dir(file)
	// 检查父目录是否存在
	if dir != "." && !IsExist(dir) {
		if err := MkdirAll(dir); err != nil {
			return nil, err
		}
	}
	return os.Create(file)
}

// IsExist
// param file string
// return bool
// 判断一个文件或目录是否存在
func IsExist(file string) bool {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		return false
	} else if err != nil {
		return false
	}
	return true
}
