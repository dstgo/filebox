package test

import (
	"fmt"
	"github.com/dstgo/filebox"
	"os"
	"testing"
)

func TestCreateFile(t *testing.T) {
	tempFile := "./testdata/dir4/createFileTest.txt"
	// 调用CreateFile方法创建文件
	file, err := filebox.CreateFile(tempFile)
	if err != nil {
		t.Errorf("Failed to create file: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	// 检查文件是否存在
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Errorf("File was not created")
	}
}
func TestClearFile(t *testing.T) {
	path := "./testdata/file1.txt"
	tempFile, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0666)

	// 向测试文件写入内容
	content := []byte("This is a test file")
	_, err = tempFile.Write(content)
	if err != nil {
		t.Errorf("Failed to write to temp file: %s", err)
	}

	// 调用ClearFile方法清空文件
	err = filebox.ClearFile(tempFile.Name())
	if err != nil {
		t.Errorf("Failed to clear file: %s", err)
	}

	// 读取文件内容
	fileContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Errorf("Failed to read file: %s", err)
	}

	// 检查文件内容是否为空
	if len(fileContent) != 0 {
		t.Errorf("File was not cleared")
	}
}
