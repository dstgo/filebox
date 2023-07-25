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

type WalkInfo struct {
	SrcPath      string
	WalkPath     string
	SrcFileInfo  fs.FileInfo
	WalkFileInfo fs.FileInfo
	Err          error
}

var (
	NotRegularFile = errors.New("not regular file")
)

// Zip
// param src string
// param dest string
// return error
// 以zip格式压缩文件，压缩文件时，仅保留输入路径的相对文件结构
func Zip(src, dest string) error {
	return ZipWith(src, dest, RelZipWalker(), ZipHeader())
}

// AppendToZip
// param zipPath string
// param sources ...string
// return error
// 将额外的文件或目录添加到现有的压缩文件中，
func AppendToZip(zipPath string, sources ...string) error {
	return AppendToZipWith(OutLayerWalker(), ZipHeader(), zipPath, sources...)
}

// Unzip
// param src string
// param dest string
// return error
// 解压zip文件到指定位置
func Unzip(src, dest string) error {
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	return unzip(zipReader, dest)
}

// Ziper 决定了用何种方式将源文件的内容写入到压缩文件中
type Ziper func(zipWriter *zip.Writer, fileHeader *zip.FileHeader, fileReader io.Reader, fileInfo fs.FileInfo) error

func ZipHeader() Ziper {
	return func(zipWriter *zip.Writer, fileHeader *zip.FileHeader, fileReader io.Reader, fileInfo fs.FileInfo) error {
		// 写入压缩文件头
		zipFile, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			return err
		}
		// 如果是文件的话就copy
		if !fileInfo.IsDir() {
			if _, err := io.CopyBuffer(zipFile, fileReader, DefaultBuffer); err != nil {
				return err
			}
		}
		return nil
	}
}

// ZipWalker 决定了用何种方式如何处理压缩信息
// 压缩时，会对输入路径执行 filepath.Walk，walkerFn会在每一次walk时调用
// 通过walker可以自定义压缩方式
type ZipWalker func(walkInfo WalkInfo, zipWriter *zip.Writer) (*zip.FileHeader, error)

// RelZipWalker 只会保留相对SrcPath的文件结构，用于直接压缩文件或目录
func RelZipWalker() ZipWalker {
	return func(walkInfo WalkInfo, zipWriter *zip.Writer) (*zip.FileHeader, error) {
		if walkInfo.Err != nil {
			return nil, walkInfo.Err
		}

		// 创建压缩文件信息头
		fileHeader, err := zip.FileInfoHeader(walkInfo.WalkFileInfo)
		if err != nil {
			return nil, err
		}
		// 指定压缩算法
		fileHeader.Method = zip.Deflate

		rel, err := filepath.Rel(walkInfo.SrcPath, walkInfo.WalkPath)
		if err != nil {
			return nil, err
		}
		// 如果得到的相对路径是"."，则说明可以直接使用文件名
		if rel == "." && !walkInfo.WalkFileInfo.IsDir() {
			rel = walkInfo.WalkFileInfo.Name()
		}
		fileHeader.Name = rel

		if !walkInfo.WalkFileInfo.IsDir() {
			if !walkInfo.WalkFileInfo.Mode().IsRegular() {
				return nil, NotRegularFile
			}
		} else if fileHeader.Name == "." {
			fileHeader.Name = "/"
		} else {
			fileHeader.Name += "/"
		}

		return fileHeader, nil
	}
}

// OutLayerWalker 也只会保留相对SrcPath的文件结构，但是相较于 RelZipWalker，
// 它还会尽量保留 path.Base(walkInfo.SrcPath)返回的最内层目录
// 用于向已存在的压缩文件添加新的文件或目录
func OutLayerWalker() ZipWalker {
	return func(walkInfo WalkInfo, zipWriter *zip.Writer) (*zip.FileHeader, error) {
		if walkInfo.Err != nil {
			return nil, walkInfo.Err
		}

		if !walkInfo.WalkFileInfo.IsDir() && walkInfo.WalkFileInfo.Mode().IsRegular() {
			return nil, NotRegularFile
		}

		// 创建压缩文件信息头
		fileHeader, err := zip.FileInfoHeader(walkInfo.WalkFileInfo)
		if err != nil {
			return nil, err
		}
		// 指定压缩算法
		fileHeader.Method = zip.Deflate

		// 如果append的是一个目录
		if walkInfo.SrcFileInfo.IsDir() {
			base := path.Base(walkInfo.SrcPath)
			if base == "." || base == "" {
				base = "/"
			}
			rel, err := filepath.Rel(walkInfo.SrcPath, walkInfo.WalkPath)
			if err != nil {
				return nil, err
			}
			fileHeader.Name = path.Join(base, rel)
		} else { // 如果是一个文件的话直接将header设置为文件名
			_, file := path.Split(walkInfo.WalkPath)
			fileHeader.Name = file
		}

		return fileHeader, nil
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
	return zipWalk(writer, src, walker, ziper)
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
		err = ziper(zipWriter, fileHeader, fileReader, zipFile.FileInfo())
		fileReader.Close()
		if err != nil {
			return err
		}
	}

	// 将新的待压缩文件写入临时zip中
	for _, src := range sources {
		if err := zipWalk(zipWriter, src, walker, ziper); err != nil {
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

// zipWalk
// param writer *zip.Writer
// param src string 目标路径
// param rel bool 如果为true，压缩时，仅保留相对文件结构
// return error
// 压缩指定路径的文件或目录创建到zip压缩文件
func zipWalk(writer *zip.Writer, src string, walker ZipWalker, ziper Ziper) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	return filepath.Walk(src, func(walkPath string, info fs.FileInfo, err error) error {
		fileHeader, err := walker(WalkInfo{
			SrcPath:      src,
			WalkPath:     walkPath,
			SrcFileInfo:  stat,
			WalkFileInfo: info,
			Err:          err,
		}, writer)
		if err != nil {
			return err
		}

		fileReader, err := OpenFileReader(walkPath)
		if err != nil {
			return err
		}
		defer fileReader.Close()

		return ziper(writer, fileHeader, fileReader, info)
	})
}

func unzip(zipReader *zip.ReadCloser, dest string) error {
	for _, zipFile := range zipReader.File {
		unzipPath := path.Join(dest, zipFile.Name)
		if !zipFile.FileInfo().IsDir() {
			// 打开压缩文件，准备读取
			reader, err := zipFile.Open()
			if err != nil {
				return err
			}

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
		} else {
			if err := MkdirAll(unzipPath); err != nil {
				return err
			}
		}
	}
	return nil
}
