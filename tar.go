package filebox

import (
	"archive/tar"
	"compress/gzip"
)

// Tar
// param src string
// param dest string
// return error
// 以tar格式压缩文件
func Tar(src, dest string) error {
	tarFile, err := CreateFile(dest)
	if err != nil {
		return err
	}
	gzWriter := gzip.NewWriter(tarFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return nil
}

// UnTar
// param src string
// param dest string
// return error
// 解压tar文件
func UnTar(src, dest string) error {
	return nil
}
