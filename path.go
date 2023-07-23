package filebox

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// GetCurrentRunningPath
// return string
// 获取当前程序的运行的绝对路径
func GetCurrentRunningPath() string {
	lookPath, _ := exec.LookPath(os.Args[0])
	abs, err := filepath.Abs(lookPath)
	if err != nil {
		return ""
	}
	return filepath.Dir(abs)
}

// GetCurrentCallerPath
// return string
// 获取调用该函数的caller的路径
func GetCurrentCallerPath() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return file
}
