package filebox

import (
	"io"
	"io/fs"
	"os"
	"path"
)

type RwFile interface {
	fs.File
	io.ReadWriteCloser
}

type RwFs interface {
	ReadFs
	WriteFs
}

type ReadFs interface {
	Open(name string) (fs.File, error)
	ReadDir(name string) ([]fs.DirEntry, error)
	ReadFile(name string) ([]byte, error)
}

type WriteFs interface {
	Mkdir(name string) error
	MkdirAll(name string) error
	WriteFile(name string, data []byte) error
	CreateFile(name string) (RwFile, error)
}

var Os = new(OsFs)

type OsFs struct {
}

func (s *OsFs) CreateFile(name string) (RwFile, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
}

func (s *OsFs) Mkdir(name string) error {
	return os.Mkdir(name, os.ModeDir)
}

func (s *OsFs) MkdirAll(name string) error {
	return os.MkdirAll(name, os.ModeDir)
}

func (s *OsFs) WriteFile(name string, data []byte) error {
	return os.WriteFile(name, data, os.ModePerm)
}

func (s *OsFs) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (s *OsFs) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (s *OsFs) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func CopyFsDir(srcFS ReadFs, dstFs WriteFs, src, dst string) error {
	dir, err := srcFS.ReadDir(src)
	if err != nil {
		return err
	} else if err = dstFs.MkdirAll(dst); err != nil {
		return err
	}

	for _, entry := range dir {
		dstPath := path.Join(dst, entry.Name())
		srcPath := path.Join(src, entry.Name())
		var copyErr error
		if entry.IsDir() {
			copyErr = CopyFsDir(srcFS, dstFs, srcPath, dstPath)
		} else {
			copyErr = CopyFsFile(srcFS, dstFs, srcPath, dstPath)
		}

		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

func CopyFsFile(srcFs ReadFs, dstFs WriteFs, src, dst string) error {
	srcFile, err := srcFs.Open(src)
	if err != nil {
		return err
	}
	dstFile, err := dstFs.CreateFile(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// CopyDir
// param src string 源路径
// param dst string 目标路径
// return error
// 将源路径的目录复制到目标路径
func CopyDir(src, dst string) error {
	return CopyFsDir(Os, Os, src, dst)
}

// CopyFile
// param src string
// param dst string
// return error
// 将源路径的文件复制到目标路径
func CopyFile(src, dst string) error {
	return CopyFsFile(Os, Os, src, dst)
}
