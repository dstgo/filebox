package filebox

import (
	"os"
	"path"
)

// Mkdir
// param dirs ...string
// return error
// 创建一个或多个目录，如果父目录不存在则返回错误
func Mkdir(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.Mkdir(dir, os.ModeDir); err != nil {
			return err
		}
	}
	return nil
}

// MkdirAll
// param dirs ...string
// return error
// 创建一个或多个目录，会尽可能的创建所有目录，包括父目录
func MkdirAll(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModeDir); err != nil {
			return err
		}
	}
	return nil
}

// MkdirTemp
// param dir string
// param pattern string
// return string
// return func() error
// return error
// 创建一个临时文件，并返回一个函数以供删除这个临时文件
func MkdirTemp(dir string, pattern string) (string, func() error, error) {
	tempDir, err := os.MkdirTemp(dir, pattern)
	rm := func() error {
		if err != nil {
			return err
		}
		return os.RemoveAll(tempDir)
	}
	return tempDir, rm, err
}

// ReadDirFullNames
// param dir string
// return []string
// 该函数返回一个路径下的所有文件或目录的完整路径
// 如果发生错误，就返回空切片
func ReadDirFullNames(dir string) []string {
	var names []string
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return names
	}
	for _, entry := range dirs {
		names = append(names, path.Join(dir, entry.Name()))
	}
	return names
}

// ReadDirShortNames
// param dir string
// return []string
// 该函数返回一个路径下的所有文件或目录的名称
// 如果发生错误，就返回空切片
func ReadDirShortNames(dir string) []string {
	var names []string
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return names
	}
	for _, entry := range dirs {
		names = append(names, entry.Name())
	}
	return names
}
