package test

import (
	"fmt"
	"github.com/dstgo/filebox"
	"testing"
)

func TestListFileNames(t *testing.T) {
	fmt.Println("------------------AllFile")
	fileNames1 := filebox.ListFileNames("./testdata", filebox.FileTypeAll)
	for _, fileName := range fileNames1 {
		fmt.Println(fileName)
	}
	fmt.Println("------------------Folder")
	fileNames2 := filebox.ListFileNames("./testdata", filebox.FileTypeFolder)
	for _, fileName := range fileNames2 {
		fmt.Println(fileName)
	}
	fmt.Println("------------------PlainFile")
	fileNames3 := filebox.ListFileNames("./testdata", filebox.FileTypePlainFile)
	for _, fileName := range fileNames3 {
		fmt.Println(fileName)
	}
	fmt.Println("------------------OtherFile")
	fileNames4 := filebox.ListFileNames("./testdata", filebox.FileTypeOther)
	for _, fileName := range fileNames4 {
		fmt.Println(fileName)
	}

}
