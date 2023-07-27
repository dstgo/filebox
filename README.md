# filebox

一个简单的文件操作工具库，封装了常用的文件操作。



## 安装

```sh
$ go get github.com/dstgo/filebox
```



## 文档

<br/>

### 状态信息

- **Stat** 返回一个文件或目录的描述信息

    ```go
    func Stat(name string) os.FileInfo
    ```

- **LStat**  返回一个符号链接的描述信息

    ```go
    func LStat(name string) os.FileInfo
    ```

- **IsExist** 判断一个文件或目录是否存在

    ```go
    func IsExist(file string) bool
    ```

- **IsLink** 判断一个是否是一个符号连接

    ```go
    func IsLink(file string) bool
    ```

- **IsDir** 判断是否是目录

    ```go
    func IsDir(path string) bool
    ```

- **IsRegular** 判断是否为标准文件

    ```go
    func IsRegular(name string) bool
    ```

- **Size** 获取一个文件的大小，单位为字节

    ```go
      func Size(name string) int64
    ```

- **FileMode** 获取文件模式位

    ```go
    func FileMode(name string) os.FileMode
    ```

- **Perm** 获取文件的权限位

    ```go
    func Perm(name string) os.FileMode
    ```

- **MTime** 获取一个文件或目录的最后修改时间

    ```go
    func MTime(file string) time.Time
    ```

### 文件操作

- **CreateFile** 检查文件的父目录是否存在，并创建文件

    ```go
    func CreateFile(file string) (*os.File, error)
    ```

- **CopyDir** 复制源目录到目标目录，如果存在则覆盖

    ```go
    func CopyDir(src, dst string) error
    ```

- **CopyFile** 复制源文件到目标文件，如果存在则覆盖

    ```go
    func CopyFile(src, dst string) error
    ```

- **CreateTempFile** 创建一个临时文件，并返回一个函数以删除这个临时文件

    ```go
    func CreateTempFile(dir, pattern string) (file *os.File, rm func() error, err error)
    ```

- **ClearFile** 清空文件内容

    ```go
    func ClearFile(path string) error 
    ```

- **ReadFileBytes**  将文件内容读取为字节切片

    ```go
    func ReadFileBytes(file string) ([]byte, error)
    ```

- **ReadFileString** 将文件内容读取为字符串

    ```go
    func ReadFileString(file string) (string, error)
    ```

- **ReadFileLine** 按行读取文件内容，返回迭代器

    ```go
    func ReadFileLine(file string) (NextLine, error) 
    ```

- **ReadFileLines** 按行读取文件内容，收集成一个字符串切片

    ```go
    func ReadFileLines(file string) ([]string, error)
    ```

### 目录操作

- **Mkdir** 创建一个或多个目录，不会检查父目录

    ```go
    func Mkdir(dirs ...string) error 
    ```

- **MkdirAll** 创建一个或多个目录，会检查父目录

    ```go
    func MkdirAll(dirs ...string) error
    ```

- **MkdirTemp** 创建一个临时文件，并返回一个函数以删除这个临时文件

    ```go
    func MkdirTemp(dir string, pattern string) (string, func() error, error)
    ```

- **ReadDirFullNames** 读取一个路径下的所有文件和目录，返回它们的完整路径

    ```go
    func ReadDirFullNames(dir string) []string
    ```

- **ReadDirShortNames** 读取一个路径下的所有文件和目录，返回它们的简短名称

    ```go
    func ReadDirShortNames(dir string) []string
    ```

- **ListFileNames** 返回指定目录下指定类型的文件名

    ```go
    func ListFileNames(dirPath string, fileType FileType) []string
    ```

### 压缩操作

* **Zip** 以zip格式压缩指定文件或目录

    ```go
    func Zip(src, dest string) error
    ```

* **AppendToZip** 将额外的文件或目录添加到已有的zip压缩文件中

    ```go
    func AppendToZip(zipPath string, sources ...string) error
    ```

* **UnZip** 解压zip压缩格式的压缩包

    ```go
    func Unzip(src, dest string) error
    ```

* **TarGzip** 以tar格式压缩指定文件或目录

    ```go
    func Tar(src, dest string) error
    ```

* **AppendToTarGzip** 向已存在的Tgz压缩文件添加新的文件或目录
  ```go
  func AppendToTarGzip(tgz string, sources ...string) error
  ```
* **UnTarGzip** 解压tar压缩格式的压缩包

    ```go
    func UnTar(src, dest string) error
    ```

### 路径操作

- **GetCurrentRunningPath** 获取当前程序运行的绝对路径

    ```go
    func GetCurrentRunningPath() string
    ```

- **GetCurrentCallerPath** 获取当前调用者的文件路径

    ```go
    func GetCurrentCallerPath() string
    ```


### 哈希操作

- **FileSum** 将一个指定路径的文件根据指定的哈希算法计算出对应的哈希值

    ```go
    func FileSum(name string, h hash.Hash) []byte 
    ```

- **FileSumReader** 根据一个reader，计算出对应的哈希值

    ```go
    func FileSumReader(reader io.Reader, h hash.Hash) []byte
    ```

- **Md5** 计算出指定路径文件的md5值

    ```go
    func Md5(name string) []byte
    ```

- **Sha1** 计算出指定路径文件的sha1值

    ```go
    func Sha1(name string) []byte
    ```

- **Sha256** 计算出指定路径文件的sha256值

    ```go
    func Sha256(name string) []byte
    ```

- **Sha512** 计算出指定路径文件的sha512值

    ```go
    func Sha512(name string) []byte
    ```

    

