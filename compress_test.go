package filebox

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestZip_UnZip(t *testing.T) {
	temp, cleanup, err := TempDir("zip")
	assert.Nil(t, err)
	if err == nil {
		defer cleanup()
	}
	file, err := CreateFile(filepath.Join(temp, "test", "test.txt"))
	assert.Nil(t, err)
	file.WriteString("hello world")
	file.Close()

	zippath := filepath.Join(os.TempDir(), "test.zip")
	err = Zip(temp, zippath)
	assert.Nil(t, err)

	err = UnZip(zippath, temp+"1")
	assert.Nil(t, err)

	content, err := os.ReadFile(filepath.Join(temp+"1", "test", "test.txt"))
	assert.Nil(t, err)
	assert.EqualValues(t, "hello world", content)
}

func TestTar_UnTar(t *testing.T) {
	temp, cleanup, err := TempDir("tar")
	assert.Nil(t, err)
	if err == nil {
		defer cleanup()
	}
	file, err := CreateFile(filepath.Join(temp, "test", "test.txt"))
	assert.Nil(t, err)
	file.WriteString("hello world")
	file.Close()

	tarpath := filepath.Join(os.TempDir(), "test.tar.gz")
	err = TarGzip(temp, tarpath)
	assert.Nil(t, err)

	err = UnTarGzip(tarpath, temp+"1")
	assert.Nil(t, err)

	content, err := os.ReadFile(filepath.Join(temp+"1", "test", "test.txt"))
	assert.Nil(t, err)
	assert.EqualValues(t, "hello world", content)
}
