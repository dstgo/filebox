package filebox

import (
	"errors"
	"io/fs"
	"os"
	"time"
)

type errorFileInfo struct {
	Error error
}

func (f errorFileInfo) Name() string {
	return ""
}

func (f errorFileInfo) Size() int64 {
	return 0
}

func (f errorFileInfo) Mode() fs.FileMode {
	return 0
}

func (f errorFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (f errorFileInfo) IsDir() bool {
	return false
}

func (f errorFileInfo) Sys() any {
	return nil
}

// Stat
// param name string
// return os.FileInfo
// 返回一个文件或目录的描述信息
func Stat(name string) os.FileInfo {
	stat, err := os.Stat(name)
	if err != nil {
		return errorFileInfo{err}
	}
	return stat
}

// LStat
// param name string
// return os.FileInfo
// 返回一个符号链接的描述信息
func LStat(name string) os.FileInfo {
	stat, err := os.Lstat(name)
	if err != nil {
		return errorFileInfo{err}
	}
	return stat
}

// ErrFromStat
// param fileInfo os.FileInfo
// return error
// 从 errorFileInfo 中获取 error
func ErrFromStat(fileInfo os.FileInfo) error {
	if info, ok := fileInfo.(errorFileInfo); ok {
		return info.Error
	}
	return nil
}

// IsDir 判断是否是目录
func IsDir(name string) bool {
	return Stat(name).IsDir()
}

// IsExist
// param name string
// return bool
// 判断指定路径上的文件或目录是否存在
func IsExist(name string) bool {
	return errors.Is(ErrFromStat(Stat(name)), os.ErrNotExist)
}

// IsLink
// param name string
// return bool
// 判断是否为符号链接
func IsLink(name string) bool {
	return LStat(name).Mode()&os.ModeSymlink != 0
}

// Size
// param name string
// return int64
// 获取文件的字节大小
func Size(name string) int64 {
	if !IsRegular(name) {
		return 0
	}
	return Stat(name).Size()
}

// IsRegular
// param name string
// return bool
// 判断是否为标准文件
func IsRegular(name string) bool {
	return Stat(name).Mode().IsRegular()
}

// FileMode 获取文件模式位
// param name string
// return os.FileMode
func FileMode(name string) os.FileMode {
	return Stat(name).Mode()
}

// MTime 获取最后的修改时间
// param name string
// return time.Time
func MTime(name string) time.Time {
	return Stat(name).ModTime()
}

// Perm 获取文件的权限位
// param name string
// return os.FileMode
func Perm(name string) os.FileMode {
	return Stat(name).Mode().Perm()
}
