package filebox

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Mkdir creates a directory, create parent-directory first if parent-directory does not exist.
func Mkdir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

// IsExist checks if the given path exists
func IsExist(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return !os.IsNotExist(err)
}

// IsDir checks if the given path is a directory
func IsDir(name string) bool {
	stat, err := os.Stat(name)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// CreateFile creates a new file, create parent-directory first if parent-directory does not exist,
// then create the specified file.
func CreateFile(filename string) (*os.File, error) {
	return OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
}

// OpenFile open a new file, create parent-directory first if parent-directory does not exist,
// then open the specified file.
func OpenFile(filename string, flag int, perm os.FileMode) (*os.File, error) {
	dir := filepath.Dir(filename)
	if dir != "." {
		if err := Mkdir(dir); err != nil {
			return nil, err
		}
	}
	file, err := os.OpenFile(filename, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// TempDir creates a temporary directory, and returns a cleanup function to clean up temporary dir
func TempDir(dir string) (temp string, cleanup func() error, err error) {
	tmpDir := filepath.Join(os.TempDir(), dir)
	if err := Mkdir(tmpDir); err != nil {
		return "", nil, err
	}
	return tmpDir, func() error {
		return os.RemoveAll(tmpDir)
	}, nil
}

// TempFile creates a temporary file, and returns a cleanup function to clean up temporary file
func TempFile(dir, pattern string) (temp *os.File, cleanup func() error, err error) {
	tempFile, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return nil, nil, err
	}
	return tempFile, func() error {
		_ = tempFile.Close()
		return os.Remove(tempFile.Name())
	}, nil
}

// ReadLine returns an iterator that iterate over per line of the given file.
func ReadLine(file string) (next func() ([]byte, error), cleanup func() error, err error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}

	reader := bufio.NewReaderSize(fd, 8192)
	lineBuf := bytes.NewBuffer(make([]byte, 8192))

	return func() ([]byte, error) {
		defer lineBuf.Reset()
		for {
			line, prefix, err := reader.ReadLine()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return lineBuf.Bytes(), nil
				}
				return nil, err
			}
			// 将读取到的当前行内容写入缓冲
			if _, err := lineBuf.Write(line); err != nil {
				return nil, err
			}

			// if this line is too long to read once a time, keep reading until finished.
			if !prefix {
				break
			}
		}
		return lineBuf.Bytes(), nil
	}, fd.Close, nil
}

// ReadLines returns a slice of bytes of file contents split by line.
func ReadLines(file string) ([][]byte, error) {
	next, cleanup, err := ReadLine(file)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	lines := make([][]byte, 0, 128)
	for {
		line, err := next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return lines, err
		}
		lines = append(lines, line)
	}
}

// ReadStringLines return a slice of strings of file contents split by line.
func ReadStringLines(filename string) ([]string, error) {
	next, cleanup, err := ReadLine(filename)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	lines := make([]string, 0, 128)
	for {
		line, err := next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return lines, err
		}
		lines = append(lines, BytesToString(line))
	}
}

// Truncate truncates the specified file
func Truncate(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if err = file.Truncate(0); err != nil {
		return err
	}
	return nil
}
