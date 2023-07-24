package filebox

import (
	"fmt"
	"os"
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

// IsDir 判断是否是目录
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func ListDirNames(dirPath string) []string {
	return listNames(dirPath, true)
}
func ListFileNames(dirPath string) []string {
	return listNames(dirPath, false)
}
func listNames(dirPath string, IsReturnDirs bool) []string {
	fileNames := make([]string, 0)

	dir, err := os.Open(dirPath)
	if err != nil {
		return fileNames
	}
	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(dir)
	files, err := dir.Readdir(-1)
	if err != nil {
		return fileNames
	}

	for _, file := range files {
		if file.IsDir() == IsReturnDirs {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames
}
