package filebox

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
	"os"
)

func FileSum(name string, h hash.Hash) []byte {
	fd, err := os.Open(name)
	if err != nil {
		return []byte{}
	}
	defer fd.Close()
	return ReaderSum(fd, h)
}

func ReaderSum(reader io.Reader, h hash.Hash) []byte {
	all, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}
	}
	h.Write(all)
	return h.Sum(nil)
}

func Sha1(name string) []byte {
	return FileSum(name, sha1.New())
}

func Sha256(name string) []byte {
	return FileSum(name, sha256.New())
}

func Sha512(name string) []byte {
	return FileSum(name, sha512.New())
}

func Md5(name string) []byte {
	return FileSum(name, md5.New())
}
