package filieutil

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// IFile 文件工具接口
type IFile interface {
	Read() ([]byte, error)                          // 基于读取当前文件数据并返回一个[]byte
	GetPath() string                                // 获取当前文件路径
	GetName() string                                // 获取当前文件名
	GetPrefix() string                              // 获取当前文件名前缀 例: test.abc 返回 test
	GetSuffix() string                              // 获取当前文件后缀 例: test.abc 返回 .abc
	Size() int64                                    // 获取当前文件大小
	IsDir() bool                                    // 获取当前文件是否是文件夹
	MkdirAll(name string) (*File, error)            // 基于当前目录创建文件夹
	Remove() error                                  // 删除当前文件
	Rename(name string) error                       // 文件重命名
	Move(path string) error                         // 文件移动
	Paste(newpath string) error                     // 文件粘贴 将当前文件粘贴到指定目录
	copyFile(src, dest string) (w int64, err error) // 内置文件复制方法
	pathExists(path string) (bool, error)           // 判断文件路径是否存在
	GetChildren() ([]*File, error)                  // 获取当前文件目录下文件数据
	Clean()                                         // 情况当前结构数据
	Replace(newFile *File)                          // 替换 将当前文件数据替换
}

type File struct {
	path  string
	isDir bool
	file  fs.FileInfo
	clean bool
}

// New 读取加载File - 基于地址获取文件结构体(整段文件逻辑入口)
func New(path string) (*File, error) {
	path = strings.ReplaceAll(path, "\\", "/")
	file, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &File{
		path:  path,
		isDir: file.IsDir(),
		file:  file,
		clean: false,
	}, nil
}

// Read 读文件数据
func (f *File) Read() ([]byte, error) {
	if f.clean {
		return nil, errors.New("it has been reset and cannot be used")
	}
	if f.isDir {
		return nil, errors.New("no such file")
	}
	return os.ReadFile(f.path)
}

// GetPath 获取文件地址
func (f *File) GetPath() string {
	return f.path
}

// GetName 获取文件名称
func (f *File) GetName() string {
	if f.clean {
		return ""
	}
	return f.file.Name()
}

// GetPrefix 获取文件名前缀
func (f *File) GetPrefix() string {
	if f.clean {
		return ""
	}
	fileName := f.GetName()
	result := fileName[0 : len(f.GetName())-len(f.GetSuffix())]
	if result != "" {
		return result
	}
	return fileName
}

// GetSuffix 获取文件后缀
func (f *File) GetSuffix() string {
	if f.clean {
		return ""
	}
	return path.Ext(f.file.Name())
}

// Size 获取文件大小
func (f *File) Size() int64 {
	if f.clean {
		return 0
	}
	return f.file.Size()
}

// IsDir 是否是文件夹
func (f *File) IsDir() bool {
	if f.clean {
		return false
	}
	return f.isDir
}

// MkdirAll 基于当前目录下创建文件夹(如果当前目录是文件则异常)
func (f *File) MkdirAll(name string) (*File, error) {
	if f.clean {
		return nil, errors.New("it has been reset and cannot be used")
	}
	if f.isDir {
		name = strings.ReplaceAll(name, "\\", "/")
		path := f.path + "/" + name
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
		return New(path)
	}
	return nil, fmt.Errorf("%s path not a folder", f.path)
}

// Remove 删除当前文件 - 删除完当前文件后回将File清空
func (f *File) Remove() error {
	if f.clean {
		return errors.New("it has been reset and cannot be used")
	}
	var err error
	if f.isDir {
		err = os.RemoveAll(f.path)
	} else {
		err = os.Remove(f.path)
	}

	if err != nil {
		return err
	}
	f.Clean()
	return nil
}

// Rename 文件重命名(文件重命名后当前file数据将变更为重命名后file数据)
func (f *File) Rename(name string) error {
	if f.clean {
		return errors.New("it has been reset and cannot be used")
	}
	path := f.path[0 : len(f.path)-len(f.file.Name())]
	err := os.Rename(f.path, path+"/"+name)
	if err != nil {
		return err
	}

	newFile, err := New(path + "/" + name)
	if err != nil {
		return err
	}
	f.Replace(newFile)
	return nil
}

// Move 文件移动(文件移动后当前file数据将变更为移动后file数据)
func (f *File) Move(path string) error {
	if f.clean {
		return errors.New("it has been reset and cannot be used")
	}
	path = strings.ReplaceAll(path, "\\", "/")
	err := os.Rename(f.path, path)
	if err != nil {
		return err
	}
	newFile, err := New(path)
	if err != nil {
		return err
	}
	f.Replace(newFile)
	return nil
}

// Paste 文件粘贴 将当前文件粘贴到指定目录(粘贴到指定目录后当前file不会变更)
func (f *File) Paste(newpath string) error {
	newpath = strings.ReplaceAll(newpath, "\\", "/")
	if f.clean {
		return errors.New("it has been reset and cannot be used")
	}
	// 判断当前目录是否为文件夹, 如果是文件夹 则需要处理文件夹内递归复制
	if f.isDir {
		// 判断path参数是否存在, 不存在则创建此目录 否则需要将文件夹内数据比对,存在同文件则覆盖, 否则写入
		if _, err := os.Stat(newpath); err != nil {
			if err := os.MkdirAll(newpath, os.ModePerm); err != nil {
				return err
			}
		}
		// 依据path路径获取copy后地址文件对象
		newFile, err := New(newpath)
		if err != nil {
			return err
		}
		if !newFile.isDir {
			return errors.New("no such directory")
		}

		return filepath.Walk(f.path, func(path string, info fs.FileInfo, err error) error {
			if info == nil {
				return err
			}
			if !info.IsDir() {
				path := strings.Replace(path, "\\", "/", -1)
				newPath := strings.Replace(path, f.path, newpath, -1)
				f.copyFile(path, newPath)
			}
			return nil
		})
	}
	// 非文件夹 直接将文件写出到对应path
	r, err := ioutil.ReadFile(f.path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(newpath, r, os.ModePerm)
}

// copyFile 内置文件复制方法
func (f *File) copyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()
	//分割path目录
	destSplitPathDirs := strings.Split(dest, "/")

	//检测时候存在目录
	var destSplitPath string
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b, _ := f.pathExists(destSplitPath)
			if !b {
				//创建目录
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

// pathExists 文件路径是否存在
func (*File) pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetChilren 获取下级目录文件数据
func (f *File) GetChildren() ([]*File, error) {
	if f.clean {
		return nil, errors.New("it has been reset and cannot be used")
	}
	if !f.isDir {
		return nil, fmt.Errorf("%s path not a folder", f.path)
	}

	children, err := ioutil.ReadDir(f.path)
	if err != nil {
		return nil, err
	}

	var result []*File
	for i := range children {
		dex := children[i]
		in := &File{
			path:  f.path + "/" + dex.Name(),
			isDir: dex.IsDir(),
			file:  dex,
		}
		result = append(result, in)
	}
	return result, nil
}

// Clean 清空自身File数据
func (f *File) Clean() {
	f.file = nil
	f.isDir = false
	f.path = ""
	f.clean = true
}

// Replace 替换 将新的file数据替换成当前file数据
func (f *File) Replace(newFile *File) {
	f.file = newFile.file
	f.isDir = newFile.isDir
	f.path = newFile.path
	f.clean = newFile.clean
}
