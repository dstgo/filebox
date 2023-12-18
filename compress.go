package filebox

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TarGzip use tar gzip to compress src path into an archive, only includes dir and regular files.
func TarGzip(src, dest string) error {
	return tarCompress(src, dest)
}

// UnTarGzip use tar gzip to decompress src archive into the destination path.
func UnTarGzip(src, dest string) error {
	return unTarCompress(src, dest)
}

// Zip use zip to compress src path into an archive, only includes dir and regular files
func Zip(src string, dest string) error {
	return zipCompress(src, dest)
}

// UnZip use zip to decompress src archive into destination path.
func UnZip(src string, dest string) error {
	return zipUnCompress(src, dest)
}

func tarCompress(src, dest string) error {
	_, err := os.Stat(src)
	if err != nil {
		return err
	}

	dir := filepath.Dir(dest)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipwriter := gzip.NewWriter(file)
	defer gzipwriter.Close()

	tarwriter := tar.NewWriter(gzipwriter)
	defer tarwriter.Close()

	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip un-regular file
		if !info.Mode().IsRegular() || info.Mode().IsDir() {
			return nil
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Size = info.Size()
		header.Name = strings.TrimPrefix(processPath(path), processPath(src))

		if err := tarwriter.WriteHeader(header); err != nil {
			return err
		}

		srcfile, err := os.Open(path)
		if err != nil {
			return err
		}

		defer srcfile.Close()

		if _, err := io.Copy(tarwriter, srcfile); err != nil {
			return err
		}

		return nil
	})
}

func unTarCompress(src, dst string) error {
	srctar, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srctar.Close()

	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	gzipreader, err := gzip.NewReader(srctar)
	if err != nil {
		return err
	}
	defer gzipreader.Close()

	tarreader := tar.NewReader(gzipreader)

	for {
		tarheader, err := tarreader.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return err
		}

		destpath := filepath.Join(dst, tarheader.Name)

		switch tarheader.Typeflag {
		case tar.TypeReg:
			destdir := filepath.Dir(destpath)
			if err := os.MkdirAll(destdir, 0755); err != nil {
				return err
			}
			destfile, err := os.OpenFile(destpath, os.O_RDWR|os.O_CREATE, os.FileMode(tarheader.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(destfile, tarreader); err != nil {
				destfile.Close()
				return err
			}
			destfile.Close()
		case tar.TypeDir:
			if err := os.MkdirAll(destpath, 0755); err != nil {
				return err
			}
		}
	}
}

func zipCompress(src string, dest string) error {
	_, err := os.Stat(src)
	if err != nil {
		return err
	}

	destWriter, err := CreateFile(dest)
	if err != nil {
		return err
	}
	defer destWriter.Close()

	zipWriter := zip.NewWriter(destWriter)
	defer zipWriter.Close()

	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relName := strings.TrimPrefix(processPath(path), processPath(src))

		if info.Mode().IsRegular() { // zip file
			reader, err := os.Open(path)
			if err != nil {
				return err
			}
			defer reader.Close()
			// zip header
			zipHeader, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			zipHeader.Name = relName
			writer, err := zipWriter.CreateHeader(zipHeader)
			if err != nil {
				return err
			}
			if _, err := io.Copy(writer, reader); err != nil {
				return err
			}
			return nil
		} else if info.Mode().IsDir() { // dir
			if relName == "." {
				relName = "/"
			} else {
				relName = filepath.Clean(relName) + "/"
			}
			if _, err := zipWriter.Create(relName); err != nil {
				return err
			}
			return nil
		} else {
			return nil
		}
	})
}

func zipUnCompress(src, dest string) error {
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	for _, fd := range zipReader.File {
		destpath := filepath.Join(dest, fd.Name)
		if fd.Mode().IsRegular() {
			if err := unzipFile(fd, destpath); err != nil {
				return err
			}
		} else if fd.Mode().IsDir() {
			if err := Mkdir(destpath); err != nil {
				return err
			}
		}
	}
	return nil
}

func unzipFile(zipfile *zip.File, destpath string) error {
	reader, err := zipfile.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := OpenFile(destpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, zipfile.Mode())
	if err != nil {
		return err
	}
	defer writer.Close()
	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}
	return nil
}

// process ignore the driver letter, replace right slash into left slash
// to make representation of path is consistent.
// eg.
// c:\\user\\appdata\\ => /user/appdata
func processPath(path string) string {
	if s := strings.Split(path, ":"); len(s) > 1 {
		path = s[1]
	}
	return strings.ReplaceAll(path, "\\", "/")
}
