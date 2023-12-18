package filebox

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir(t *testing.T) {
	src, cleanupSrc, err := TempDir("src")
	assert.Nil(t, err)
	if err == nil {
		defer func() {
			fmt.Println(cleanupSrc())
		}()
	}

	dest, cleanupDest, err := TempDir("dest")
	assert.Nil(t, err)
	if err == nil {
		defer func() {
			fmt.Println(cleanupDest())
		}()
	}

	file, err := CreateFile(filepath.Join(src, "hello", "hello.txt"))
	assert.Nil(t, err)
	_, err = file.Write([]byte("hello"))
	assert.Nil(t, err)
	err = file.Close()
	assert.Nil(t, err)

	err = Copy(src, dest)
	t.Log(err)
	assert.Nil(t, err)

	content, err := os.ReadFile(filepath.Join(dest, "hello", "hello.txt"))
	assert.Nil(t, err)
	assert.EqualValues(t, []byte("hello"), content)
}

func TestCopyFile(t *testing.T) {
	temp1, cleanup1, err := TempFile(os.TempDir(), "hello1")
	assert.Nil(t, err)
	if err == nil {
		defer cleanup1()
	}
	temp2, cleanup2, err := TempFile(os.TempDir(), "hello2")
	assert.Nil(t, err)
	if err == nil {
		temp2.Close()
		defer cleanup2()
	}

	_, err = temp1.Write([]byte("hello"))
	assert.Nil(t, err)

	err = Copy(temp1.Name(), temp2.Name())
	assert.Nil(t, err)

	content, err := os.ReadFile(temp2.Name())
	assert.Nil(t, err)
	assert.EqualValues(t, []byte("hello"), content)
}
