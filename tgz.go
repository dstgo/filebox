package filebox

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

// TarGzip
// param src string
// param dest string
// return error
// 以tgz格式压缩文件
// 如果默认的walker和tgz不能满足要求
// 可以使用 TarGzipLevelWith 来传入自定义的walker和tgz
func TarGzip(src, dest string) error {
	return TarGzipLevelWith(src, dest, gzip.DefaultCompression, DefaultWalker(), DefaultTgz())
}

// TarGzipLevelWith
// param src string
// param dest string
// param level int
// return error
// 以tgz格式压缩文件，可以指定压缩级别
func TarGzipLevelWith(src, dest string, level int, walker TgzWalker, tgz Tgz) error {
	tarFile, err := CreateFile(dest)
	if err != nil {
		return err
	}
	defer tarFile.Close()
	gzWriter, err := gzip.NewWriterLevel(tarFile, level)
	if err != nil {
		return err
	}
	defer gzWriter.Close()

	// 用tar.Writer包装gzip.Writer，达到gzip压缩的目的
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return tgzWalk(tarWriter, walker, tgz, src)
}

// AppendToTarGzip 将额外的文件或目录添加到已存在的tar.gz归档文件中
// param tgz string
// param sources ...string
// return error
// 如果默认的walker和tgz无法满足要求
// 可以使用 AppendToTarGzipLevelWith 来进行自定义
func AppendToTarGzip(tgz string, sources ...string) error {
	return AppendToTarGzipLevelWith(gzip.DefaultCompression, DefaultWalker(), DefaultTgz(), tgz, sources...)
}

// AppendToTarGzipLevelWith
// param level int
// param walker TgzWalker
// param tgz Tgz
// param tgzPath string
// param sources ...string
// return error
// 将额外的文件或目录添加到已存在的tar.gz归档文件中un
// 对于tar已经归档的文件，是无法像zip一样再次更新内容的
// 如果要做到类似的效果需要先解压再压缩
func AppendToTarGzipLevelWith(level int, walker TgzWalker, tgz Tgz, tgzPath string, sources ...string) error {

	// 创建一个临时文件夹以存放解压文件
	tempDir, rmDir, err := MkdirTemp(os.TempDir(), "untgz")
	if err != nil {
		return err
	}
	defer rmDir()

	// 将tgz解压到临时文件夹
	if err := UnTarGzip(tgzPath, tempDir); err != nil {
		return err
	}

	// 将这些已经解压过了的添加到待压缩列表中
	untgzItems := ReadDirFullNames(tempDir)
	sources = append(untgzItems, sources...)

	// 创建临时tgz文件
	temFileReader, rm, err := CreateTempFile(os.TempDir(), "temp.tgz")
	if err != nil {
		return err
	}
	defer rm()

	tempGzipWriter, err := gzip.NewWriterLevel(temFileReader, level)
	if err != nil {
		return err
	}

	tempTarWriter := tar.NewWriter(tempGzipWriter)

	// 添加tgz压缩
	for _, src := range sources {
		if err := tgzWalk(tempTarWriter, walker, tgz, src); err != nil {
			return err
		}
	}

	if err := errors.Join(
		tempTarWriter.Close(),
		tempGzipWriter.Close(),
		temFileReader.Close()); err != nil {
		return err
	}

	// 覆盖原有的压缩文件
	if err := CopyFile(temFileReader.Name(), tgzPath); err != nil {
		return err
	}

	return nil
}

type TgzWalker func(info WalkInfo, writer *tar.Writer) (*tar.Header, error)

func DefaultWalker() TgzWalker {
	return func(walkInfo WalkInfo, writer *tar.Writer) (*tar.Header, error) {
		if walkInfo.Err != nil {
			return nil, walkInfo.Err
		}

		var link string

		if IsLink(walkInfo.WalkPath) {
			targetLink, err := os.Readlink(walkInfo.WalkPath)
			if err != nil {
				return nil, err
			}
			link = targetLink
		}

		fileHeader, err := tar.FileInfoHeader(walkInfo.WalkFileInfo, link)
		if err != nil {
			return nil, err
		}

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

type Tgz func(tarWriter *tar.Writer, fileHeader *tar.Header, fileReader io.Reader, fileInfo fs.FileInfo) error

func DefaultTgz() Tgz {
	return func(tarWriter *tar.Writer, fileHeader *tar.Header, fileReader io.Reader, fileInfo fs.FileInfo) error {
		if err := tarWriter.WriteHeader(fileHeader); err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			_, err := io.CopyBuffer(tarWriter, fileReader, DefaultBuffer)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func tgzWalk(tarWriter *tar.Writer, walker TgzWalker, tgz Tgz, src string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	return filepath.Walk(src, func(walkPath string, info fs.FileInfo, err error) error {
		fileHeader, err := walker(WalkInfo{
			SrcPath:      src,
			WalkPath:     walkPath,
			SrcFileInfo:  srcStat,
			WalkFileInfo: info,
			Err:          err,
		}, tarWriter)

		fileReader, err := OpenFileReader(walkPath)
		if err != nil {
			return err
		}
		defer fileReader.Close()
		return tgz(tarWriter, fileHeader, fileReader, info)
	})
}

// UnTarGzip
// param src string
// param dest string
// return error
// 解压tgz文件
func UnTarGzip(src, dest string) error {
	srcReader, err := OpenFileReader(src)
	if err != nil {
		return err
	}
	defer srcReader.Close()

	gzReader, err := gzip.NewReader(srcReader)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	return unTgz(tarReader, dest)
}

func unTgz(tarReader *tar.Reader, dest string) error {
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}

		unTgzPath := path.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir: // dir
			if err := MkdirAll(unTgzPath); err != nil {
				return err
			}
		case tar.TypeReg: // regular file
			fileWriter, err := CreateFileMode(unTgzPath, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			_, err = io.CopyBuffer(fileWriter, tarReader, DefaultBuffer)
			fileWriter.Close()
			if err != nil {
				return err
			}
		case tar.TypeSymlink: // symbol link
			if err := os.Symlink(header.Linkname, unTgzPath); err != nil {
				return err
			}
		}
	}
}
