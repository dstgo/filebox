# filebox
一个简单的文件操作工具库，封装了常用的文件操作。



##  安装

```sh
$ go get github.com/dstgo/filebox
```



## 压缩操作

* **Zip** 以zip格式压缩指定文件或目录

    ```go
    func Zip(src, dest string) error
    ```
    
* **UnZip** 解压zip压缩格式的压缩包

    ```go
    func Unzip(src, dest string) error
    ```

* **Tar** 以tar格式压缩指定文件或目录

    ```go
    func Tar(src, dest string) error
    ```

* **UnTar** 解压tar压缩格式的压缩包

    ```go
    func UnTar(src, dest string) error
    ```
## 路径操作

- **GetCurrentRunningPath** 获取当前程序运行的绝对路径

    ```go
    func GetCurrentRunningPath() string
    ```

- **GetCurrentCallerPath** 获取当前调用者的文件路径

    ```go
    func GetCurrentCallerPath() string
    ```

    

## 文件操作

- **CreateFile** 检查文件的父目录是否存在，并创建文件

    ```go
    func CreateFile(file string) (*os.File, error)
    ```

- **CopyDir** 复制源目录到目标目录

    ```go
    func CopyDir(src, dst string) error
    ```

- **CopyFile** 复制源文件到目标文件

    ```go
    func CopyFile(src, dst string) error
    ```

- **ClearFile** 清空文件内容
    
  ```go
  func ClearFile(path string) error 
  ```

## 目录操作

- **Mkdir** 创建一个或多个目录，不会检查父目录

    ```go
    func Mkdir(dirs ...string) error 
    ```

- **MkdirAll** 创建一个或多个目录，会检查父目录

    ```go
    func MkdirAll(dirs ...string) error
    ```

- **IsDir** 判断是否是目录

    ```go
    func IsDir(path string) bool
    ```

- **ListFileNames** 返回指定目录下指定类型的文件名

    ```go
    func ListFileNames(dirPath string, fileType FileType) []string
    ```
    

