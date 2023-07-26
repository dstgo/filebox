package filebox

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"time"
)

var (
	WriteFlag  = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	AppendFlag = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func OpenFileReader(file string) (*os.File, error) {
	return OpenFile(file, os.O_RDONLY, 0)
}

func OpenFileWriter(file string) (*os.File, error) {
	if file, err := CreateFile(file); err != nil {
		return file, err
	}
	return OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
}

func OpenFileRw(file string) (*os.File, error) {
	if file, err := CreateFile(file); err != nil {
		return file, err
	}
	return OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

// CreateFile
// param file string
// return *os.File
// return error
// 创建一个指定名称的文件，并且会检查文件的父目录是否存在
func CreateFile(file string) (*os.File, error) {
	return CreateFileMode(file, 0666)
}

// CreateFileMode
// param file string
// param mode os.FileMode
// return *os.File
// return error
// 创建一个指定名称和mode的文件，并且会检查文件的父目录是否存在
func CreateFileMode(file string, mode os.FileMode) (*os.File, error) {
	dir := path.Dir(file)
	// 检查父目录是否存在
	if dir != "." && !IsExist(dir) {
		if err := MkdirAll(dir); err != nil {
			return nil, err
		}
	}
	return OpenFile(file, WriteFlag, mode)
}

// CreateTempFile
// param dir string
// param pattern string
// return rm func() error 删除临时文件
// return err error
// 创建一个临时文件，并返回一个函数以删除这个临时文件
func CreateTempFile(dir, pattern string) (file *os.File, rm func() error, err error) {
	tempFile, err := os.CreateTemp(dir, pattern)
	rm = func() error {
		if err == nil {
			return os.Remove(tempFile.Name())
		}
		return err
	}
	return tempFile, rm, err
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

// FileSize
// param file string
// return int64
// 获取一个文件的大小
func FileSize(file string) int64 {
	stat, err := os.Stat(file)
	if err != nil {
		return 0
	}
	if stat.Mode().IsRegular() {
		return stat.Size()
	}
	return 0
}

// MTime
// param file string
// return time.Time
// 获取指定路径的文件或目录的修改时间
func MTime(file string) time.Time {
	stat, err := os.Stat(file)
	if err != nil {
		return time.Time{}
	}
	return stat.ModTime()
}

// IsLink
// param file string
// return bool
// 判断一个指定路径的目标是不是符号链接
func IsLink(file string) bool {
	lstat, err := os.Lstat(file)
	if err != nil {
		return false
	}
	return isLink(lstat)
}

func isLink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

// ReadFileBytes
// param file string
// return []byte
// return error
// 从文件里面读取字节切片
func ReadFileBytes(file string) ([]byte, error) {
	return os.ReadFile(file)
}

// ReadFileString
// param file string
// return string
// return error
// 从文件里面读取字符串
func ReadFileString(file string) (string, error) {
	bc, err := ReadFileBytes(file)
	if err != nil {
		return "", err
	}
	return string(bc), nil

}

type NextLine = func() ([]byte, error)

// ReadFileLine
// param file string
// return []string
// return error
// 成功打开文件后，返回一个next函数，用于按行读取文件内容
func ReadFileLine(file string) (NextLine, error) {
	f, err := OpenFileReader(file)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(f)
	// 行缓冲
	bufline := bytes.NewBuffer(make([]byte, 0, 4096))

	return func() ([]byte, error) {
		// 如果当前行大于缓冲区，一次读不完，则尝试多次读取，直到读取完毕
		// defaultBufSize = 4096
		for {
			line, prefix, err := reader.ReadLine()
			if err != nil {
				f.Close()
				return bufline.Bytes(), err
			}
			// 将读取到的当前行内容写入缓冲
			if _, err := bufline.Write(line); err != nil {
				return nil, err
			}

			// 当前行读取完毕就退出循环
			if !prefix {
				break
			}
		}
		defer bufline.Reset()
		return bufline.Bytes(), nil
	}, nil
}

// ReadFileLines
// param file string
// return []string
// return error
// 一次性读完文件的所有内容，并返回字符串切片
// 如果不想一次性读完全部内容，可以使用 ReadFileLine
func ReadFileLines(file string) ([]string, error) {
	nextLine, err := ReadFileLine(file)
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0, 64)
	for {
		line, err := nextLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return lines, err
		}
		lines = append(lines, string(line))
	}
}

// ClearFile 清空一个文件
func ClearFile(path string) error {
	if IsDir(path) {
		return errors.New("the path file is not a single file")
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	if err = file.Truncate(0); err != nil {
		return err
	}
	return nil
}
