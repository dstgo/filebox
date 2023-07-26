package filebox

import (
	"os"
)

type FileType int

const (
	FileTypeAll       FileType = iota
	FileTypePlainFile          // 文件
	FileTypeFolder             // 文件夹
	FileTypeOther              // 其他类型
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

func ListFileNames(dirPath string, fileType FileType) []string {
	var fileNames []string

	dir, err := os.Open(dirPath)
	if err != nil {
		return fileNames
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return fileNames
	}

	for _, file := range files {
		switch fileType {
		case FileTypeAll:
			fileNames = append(fileNames, file.Name())
			break
		case FileTypeFolder:
			if file.IsDir() {
				fileNames = append(fileNames, file.Name())
			}
			break
		case FileTypePlainFile:
			if !file.IsDir() && file.Mode().IsRegular() {
				fileNames = append(fileNames, file.Name())
			}
			break
		case FileTypeOther:
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames
}
