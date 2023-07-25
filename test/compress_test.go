package test

import (
	"fmt"
	"github.com/dstgo/filebox"
	"testing"
)

func TestZip(t *testing.T) {
	// 压缩单个文件
	fmt.Println(filebox.Zip("/test/log.txt", "/test/log.zip"))
	// 压缩一个目录
	fmt.Println(filebox.Zip("/test/home", "/test/home.zip"))
	// 相对路径压缩单个文件
	fmt.Println(filebox.Zip("./log.txt", "./log.zip"))
	// 相对路径压缩一个目录
	fmt.Println(filebox.Zip("./log", "./logd.zip"))
}

func TestAppendZip(t *testing.T) {
	// 压缩单个文件
	fmt.Println(filebox.Zip("/test/log.txt", "/test/log.zip"))
	// 添加压缩文件
	fmt.Println(filebox.AppendToZip("/test/log.zip", "/test/bob", "/test/log1.txt"))
	// 向不存在的zip添加文件
	fmt.Println(filebox.AppendToZip("/test/logaaa.zip", "/test/bob", "/test/log1.txt"))
}

func TestUnzip(t *testing.T) {
	fmt.Println(filebox.Unzip("/test/log.zip", "/test/unzip/"))
	fmt.Println(filebox.Unzip("/test/log.zip", "./unzip/"))
}

func TestTgz(t *testing.T) {
	fmt.Println(filebox.TarGzip("/test/unzip", "/test/unzip.tar.gz"))
	fmt.Println(filebox.AppendToTarGzip("/test/unzip.tar.gz", "/test/aaa.txt"))
}

func TestUnTgz(t *testing.T) {
	fmt.Println(filebox.UnTarGzip("/test/unzip.tgz", "/test/tar/"))
}
