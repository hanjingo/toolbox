package gen

import (
	"errors"
	"os"
	"path/filepath"
)

func mustOpenFile(fileName string, flag int) (*os.File, bool, error) {
	filePathName := fileName
	path, _ := filepath.Split(filePathName)
	if !isExist(filePathName) {
		if !isExist(path) {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return nil, false, err
			}
		}
		fd, err := createFile(filePathName)
		if err != nil {
			return nil, false, err
		}
		return fd, true, nil
	} else {
		if !isFile(filePathName) {
			return nil, false, errors.New("无法读取路径信息")
		}
		var err error
		fd, err := os.OpenFile(filePathName, flag, 0666)
		if err != nil {
			return nil, false, err
		}
		return fd, false, nil
	}
}

/*获得当前程序的运行路径*/
func GetCurrPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

//判断是否是文件
func isFile(filePathName string) bool {
	if f, err := os.Stat(filePathName); err == nil {
		if !f.IsDir() {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//判断文件/路径是否存在
func isExist(arg string) bool {
	_, err := os.Stat(arg) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

/*创建文件*/
func createFile(arg string) (*os.File, error) {
	if !isExist(arg) { //如果文件不存在
		fd, err := os.Create(arg)
		if err != nil {
			return nil, err
		}
		return fd, nil
	}
	return nil, errors.New("文件已经存在")
}

/*清空文件内容*/
func cleanFile(arg string) (*os.File, error) {
	if !isExist(arg) {
		return nil, errors.New("文件不存在")
	}
	if !isFile(arg) {
		return nil, errors.New("不是文件")
	}
	if err := os.Remove(arg); err != nil {
		return nil, err
	}
	fd, err := os.OpenFile(arg, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return fd, nil
}
