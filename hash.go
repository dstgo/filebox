package filebox

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
)

// FileSum 将一个指定路径的文件根据指定的哈希算法计算出对应的哈希值
// param name string
// param h hash.Hash
// return []byte
func FileSum(name string, h hash.Hash) []byte {
	reader, err := OpenFileReader(name)
	if err != nil {
		return []byte{}
	}
	return FileSumReader(reader, h)
}

// FileSumReader 根据一个reader，计算出对应的哈希值
// param reader io.Reader
// param h hash.Hash
// return []byte
func FileSumReader(reader io.Reader, h hash.Hash) []byte {
	all, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}
	}
	h.Write(all)
	return h.Sum(nil)
}

// Sha1 计算出指定路径文件的sha1值
// param name string
// return []byte
func Sha1(name string) []byte {
	return FileSum(name, sha1.New())
}

// Sha256 计算出指定路径文件的sha256值
// param name string
// return []byte
func Sha256(name string) []byte {
	return FileSum(name, sha256.New())
}

// Sha512 计算出指定路径文件的sha512值
// param name string
// return []byte
func Sha512(name string) []byte {
	return FileSum(name, sha512.New())
}

// Md5 计算出指定路径文件的md5值
// param name string
// return []byte
func Md5(name string) []byte {
	return FileSum(name, md5.New())
}
