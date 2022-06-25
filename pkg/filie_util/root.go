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
	Read() ([]byte, error)                          // based on reading current file data and return the file bytes
	GetPath() string                                // get the file path
	GetName() string                                // get the file name
	GetPrefix() string                              // get the file prefix name demo: test.abc return test
	GetSuffix() string                              // get the file suffix name demo: test.abc return .abc
	Size() int64                                    // get the file size
	IsDir() bool                                    // the file is a folder
	MkdirAll(name string) (*File, error)            // based on the current folder create a new folder
	Remove() error                                  // remove the current file
	Rename(name string) error                       // based on the current file rename to a new file
	Move(path string) error                         // based on the current file move to a new file
	Paste(newpath string) error                     // based on the current file paste to a new file
	copyFile(src, dest string) (w int64, err error) // private function to copy file
	pathExists(path string) (bool, error)           // private function to check file path exists
	GetChildren() ([]*File, error)                  // based on the current folder get all the children files
	Clean()                                         // clean the file
	Replace(newFile *File)                          // based on the current file replace to a new file
}

type File struct {
	path  string
	isDir bool
	file  fs.FileInfo
	clean bool
}

// New read the file path and return a file object
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

// Read based on reading current file data and return the file bytes
func (f *File) Read() ([]byte, error) {
	if f.clean {
		return nil, errors.New("it has been reset and cannot be used")
	}
	if f.isDir {
		return nil, errors.New("no such file")
	}
	return os.ReadFile(f.path)
}

// GetPath get the file path
func (f *File) GetPath() string {
	return f.path
}

// GetName get the file name
func (f *File) GetName() string {
	if f.clean {
		return ""
	}
	return f.file.Name()
}

// GetPrefix get the file prefix name demo: test.abc return test
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

// GetSuffix get the file suffix name demo: test.abc return .abc
func (f *File) GetSuffix() string {
	if f.clean {
		return ""
	}
	return path.Ext(f.file.Name())
}

// Size get the file size
func (f *File) Size() int64 {
	if f.clean {
		return 0
	}
	return f.file.Size()
}

// IsDir the file is a folder
func (f *File) IsDir() bool {
	if f.clean {
		return false
	}
	return f.isDir
}

// MkdirAll based on the current folder create a new folder
// if the name path is file then return error
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

// Remove remove the current file
// remove the current file the file object will be reset
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

// Rename based on the current file rename to a new file
// file rename to a new file
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

// Move based on the current file move to a new file
// file move to a new file
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

// Paste based on the current file paste to a new file
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

// copyFile private function to copy file
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

// pathExists private function to check file path exists
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

// GetChilren based on the current folder get all the children files
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

// Clean clean the file
func (f *File) Clean() {
	f.file = nil
	f.isDir = false
	f.path = ""
	f.clean = true
}

// Replace based on the current file replace to a new file
func (f *File) Replace(newFile *File) {
	f.file = newFile.file
	f.isDir = newFile.isDir
	f.path = newFile.path
	f.clean = newFile.clean
}
