package filebox

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// BinaryPath returns the current binary running path
func BinaryPath() string {
	lookPath, _ := exec.LookPath(os.Args[0])
	abs, err := filepath.Abs(lookPath)
	if err != nil {
		return ""
	}
	return filepath.Dir(abs)
}

// CallerPath returns the caller file path
func CallerPath() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return file
}
