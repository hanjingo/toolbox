package gen

import (
	"strings"
)

//排序
func SortWithId(arg map[string]*Class) []*Class {
	var back []*Class
	for _, item := range arg {
		back = append(back, item)
	}
	if len(back) <= 1 {
		return back
	}
	for i := 0; i < len(back)-1; i++ {
		for j := i; j < len(back); j++ {
			if back[i].Id > back[j].Id {
				//swap
				temp := back[i]
				back[i] = back[j]
				back[j] = temp
			}
		}
	}
	return back
}

//格式化一行 a columeSize+1字 b 24字 c
func formatLine(size int, args ...string) string {
	back := ""
	colSize := defSize
	if size > 0 {
		colSize = size
	}
	i := 0
	for _, arg := range args {
		i++
		back += arg
		if i < len(args) {
			for j := 0; len(arg)+j < colSize; j++ {
				back += " "
			}
		}
	}
	back += " "
	return back
}

//是否打印id
func isPrintId(fileMap map[string]string) bool {
	for k, path := range fileMap {
		if strings.ToUpper(k) == KEY_ID && path != "" {
			return true
		}
	}
	return false
}

//是否打印model
func isPrintModel(fileMap map[string]string) bool {
	for k, path := range fileMap {
		if strings.ToUpper(k) == KEY_MODEL && path != "" {
			return true
		}
	}
	return false
}

//是否打印err
func isPrintErr(fileMap map[string]string) bool {
	for k, path := range fileMap {
		if strings.ToUpper(k) == KEY_ERR && path != "" {
			return true
		}
	}
	return false
}

//是否打印文档model
func isPrintDocModel(fileMap map[string]string) bool {
	for k, path := range fileMap {
		if strings.ToUpper(k) == KEY_DOC_MODEL && path != "" {
			return true
		}
	}
	return false
}

//是否打印文档err
func isPrintDocErr(fileMap map[string]string) bool {
	for k, path := range fileMap {
		if strings.ToUpper(k) == KEY_DOC_ERR && path != "" {
			return true
		}
	}
	return false
}

//检查结果
func check(f func() error) string {
	err := f()
	if err != nil {
		return err.Error()
	}
	return "成功"
}
