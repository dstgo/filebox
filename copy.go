package filebox

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Copy copies src to dest path, only copy dir and regular files.
func Copy(src, dst string) error {
	return CopyFs(os.DirFS(src), ".", dst)
}

// CopyFs copies the specified srcpath from srcFs to destpath in os fs, only copy dir and regular files.
func CopyFs(srcFs fs.FS, srcpath, destpath string) error {
	// check src stat
	_, err := fs.Stat(srcFs, srcpath)
	if err != nil {
		return err
	}

	return fs.WalkDir(srcFs, srcpath, func(name string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// destination path
		destName := filepath.Join(destpath, name)

		if info.Type().IsRegular() {
			fileInfo, err := info.Info()
			if err != nil {
				return err
			}
			destfile, err := OpenFile(destName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileInfo.Mode())
			if err != nil {
				return err
			}
			defer destfile.Close()
			srcfile, err := srcFs.Open(name)
			if err != nil {
				return err
			}
			defer srcfile.Close()
			// copy file contents
			if _, err := io.Copy(destfile, srcfile); err != nil {
				return err
			}
			return nil
		} else if info.Type().IsRegular() {
			return Mkdir(destName)
		} else {
			// ignore other types
			return nil
		}
	})
}
