package test

import (
	"github.com/dstgo/filebox"
	"reflect"
	"testing"
)

func TestListDirNames(t *testing.T) {
	dirPath := "./testdata"

	// 调用ListDirNames函数获取目录名称列表
	dirNames := filebox.ListDirNames(dirPath)

	// 检查返回的目录名称列表是否符合预期
	expectedDirNames := []string{"dir1", "dir2", "dir3"}
	if !reflect.DeepEqual(dirNames, expectedDirNames) {
		t.Errorf("ListDirNames(%q) returned %v, expected %v", dirPath, dirNames, expectedDirNames)
	}
}

func TestListFileNames(t *testing.T) {
	dirPath := "./testdata"

	// 调用ListFileNames函数获取非目录文件名称列表
	fileNames := filebox.ListFileNames(dirPath)

	// 检查返回的非目录文件名称列表是否符合预期
	expectedFileNames := []string{"file1.txt", "file2.md"}
	if !reflect.DeepEqual(fileNames, expectedFileNames) {
		t.Errorf("ListFileNames(%q) returned %v, expected %v", dirPath, fileNames, expectedFileNames)
	}
}
