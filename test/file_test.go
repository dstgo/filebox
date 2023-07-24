package test

import (
	"fmt"
	"github.com/dstgo/filebox"
	"testing"
)

func TestCreateFile(t *testing.T) {
	fmt.Println(filebox.CreateFile("/test/log.txt"))
	fmt.Println(filebox.CreateFile("log.txt"))
	fmt.Println(filebox.CreateFile("./log.txt"))
}

func TestExist(t *testing.T) {
	fmt.Println(filebox.IsExist("/test/log.txt"))
}

func TestReadLines(t *testing.T) {
	fmt.Println(filebox.ReadFileLines("/test/log.txt"))
}
