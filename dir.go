package filebox

import "os"

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
