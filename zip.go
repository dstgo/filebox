package filebox

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

// Zip
// param src string
// param dest string
// return error
// 以zip格式压缩文件，压缩文件时，仅保留输入路径的相对文件结构
func Zip(src, dest string) error {
	return ZipWith(src, dest, RelZipWalker(true), ZipHeader())
}

// AppendToZip
// param zipPath string
// param sources ...string
// return error
// 将额外的文件或目录添加到现有的压缩文件中，
func AppendToZip(zipPath string, sources ...string) error {
	return AppendToZipWith(RelZipWalker(false), ZipHeader(), zipPath, sources...)
}

type ZipWalkInfo struct {
	Src      string
	WalkPath string
	FileInfo fs.FileInfo
	Err      error
}

// Ziper 决定了用何种方式将源文件的内容写入到压缩文件中
type Ziper func(zipWriter *zip.Writer, fileHeader *zip.FileHeader, fileReader io.Reader) error

func ZipHeader() Ziper {
	return func(zipWriter *zip.Writer, fileHeader *zip.FileHeader, fileReader io.Reader) error {
		zipFile, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			return err
		}
		if _, err := io.CopyBuffer(zipFile, fileReader, DefaultBuffer); err != nil {
			return err
		}
		return nil
	}
}

// ZipWalker 决定了用何种方式去遍历源文件以及如何处理压缩信息
// 压缩时，会对输入路径执行 filepath.Walk，walkerFn会在每一次walk时调用
// 通过walker可以自定义压缩方式
type ZipWalker func(walkInfo *ZipWalkInfo, zipWriter *zip.Writer, zipFn Ziper) error

func RelZipWalker(rel bool) ZipWalker {
	return func(walkInfo *ZipWalkInfo, zipWriter *zip.Writer, ziper Ziper) error {
		if walkInfo.Err != nil {
			return walkInfo.Err
		}
		if !walkInfo.FileInfo.IsDir() {
			// 创建压缩文件信息头
			fileHeader, err := zip.FileInfoHeader(walkInfo.FileInfo)
			if err != nil {
				return err
			}

			fileHeader.Name = walkInfo.WalkPath
			if rel {
				rel, err := filepath.Rel(walkInfo.Src, walkInfo.WalkPath)
				if err != nil {
					return err
				}
				// 如果得到的相对路径是"."，则说明可以直接使用文件名
				if rel == "." {
					rel = walkInfo.FileInfo.Name()
				}
				fileHeader.Name = rel
			}
			// 指定压缩算法
			fileHeader.Method = zip.Deflate

			srcFileReader, err := OpenFileReader(walkInfo.WalkPath)
			if err != nil {
				return err
			}
			defer srcFileReader.Close()

			// 压缩文件
			if err := ziper(zipWriter, fileHeader, srcFileReader); err != nil {
				return err
			}
		}
		return nil
	}
}

// ZipWith
// param src string
// param dest string
// param walker ZipWalker
// return error
// 使用walker压缩文件
func ZipWith(src, dest string, walker ZipWalker, ziper Ziper) error {
	zipFile, err := CreateFile(dest)
	defer zipFile.Close()
	if err != nil {
		return err
	}
	writer := zip.NewWriter(zipFile)
	defer writer.Close()
	return zipArchive(writer, src, walker, ziper)
}

// AppendToZipWith
// param zipPath string zip压缩文件路径
// param sources ...string 待添加的源文件路径
// return error
// 添加一个或多个文件或目录到指定的zip压缩文件中
func AppendToZipWith(walker ZipWalker, ziper Ziper, zipPath string, sources ...string) error {
	// 创建临时的zip文件
	temZipFile, rm, err := CreateTempFile(os.TempDir(), "temp.zip")
	if err != nil {
		return err
	}

	// 打开原有的zip
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}

	// 打开临时zip
	zipWriter := zip.NewWriter(temZipFile)

	// 将原有的zip写入临时zip中
	for _, zipFile := range zipReader.File {
		fileHeader, err := zip.FileInfoHeader(zipFile.FileInfo())
		if err != nil {
			return err
		}
		fileHeader.Name = zipFile.Name

		fileReader, err := zipFile.Open()
		if err != nil {
			return err
		}
		err = ziper(zipWriter, fileHeader, fileReader)
		fileReader.Close()
		if err != nil {
			return err
		}
	}

	// 将新的待压缩文件写入临时zip中
	for _, src := range sources {
		if err := zipArchive(zipWriter, src, walker, ziper); err != nil {
			return err
		}
	}

	// 关闭这几个文件
	if err := errors.Join(
		zipReader.Close(),
		// 必须先关闭zipWriter，再关闭tempfile
		zipWriter.Close(),
		temZipFile.Close(),
	); err != nil {
		return err
	}

	// 将临时的zip覆盖原有的zip
	if err := CopyFile(temZipFile.Name(), zipPath); err != nil {
		return err
	}

	// 删除临时文件
	return rm()
}

// zipArchive
// param writer *zip.Writer
// param src string 目标路径
// param rel bool 如果为true，压缩时，仅保留相对文件结构
// return error
// 压缩指定路径的文件或目录创建到zip压缩文件
func zipArchive(writer *zip.Writer, src string, walker ZipWalker, ziper Ziper) error {
	return filepath.Walk(src, func(walkPath string, info fs.FileInfo, err error) error {
		return walker(&ZipWalkInfo{
			Src:      src,
			WalkPath: walkPath,
			FileInfo: info,
			Err:      err,
		}, writer, ziper)
	})
}

// Unzip
// param src string
// param dest string
// return error
// 解压zip文件
func Unzip(src, dest string) error {
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	return unzip(zipReader, dest)
}

func unzip(zipReader *zip.ReadCloser, dest string) error {
	for _, zipFile := range zipReader.File {
		if !zipFile.FileInfo().IsDir() {
			// 打开压缩文件，准备读取
			reader, err := zipFile.Open()
			if err != nil {
				return err
			}

			unzipPath := path.Join(dest, zipFile.Name)

			if err := MkdirAll(path.Dir(unzipPath)); err != nil {
				return err
			}

			// 打开目标文件，准备写入
			writer, err := OpenFile(unzipPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, zipFile.Mode())
			if err != nil {
				return err
			}

			_, err = io.CopyBuffer(writer, reader, DefaultBuffer)

			if err := errors.Join(err, writer.Close(), reader.Close()); err != nil {
				return err
			}
		}
	}
	return nil
}
