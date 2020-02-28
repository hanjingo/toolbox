package gen

import (
	"io/ioutil"
	"os"
	"strings"
)

type GoGenerator1 struct {
	Items map[string]*Class //key:msg name  value:item集合
}

func NewGoGenerator1(conf *Config) *GoGenerator1 {
	back := &GoGenerator1{
		Items: make(map[string]*Class),
	}
	start := 0
	for _, e := range conf.Classes {
		if e.Id != 0 {
			start = e.Id
		} else {
			start++
		}
		e.Id = start
		class := NewClass(e)
		back.Items[class.Name] = class
	}
	return back
}

// msgid.go
func (gen *GoGenerator1) GenMsgid() error {
	//msgid

	items := SortWithId(gen.Items)
	for _, item := range items {
		var namespace = MSGID_PACK_NAME
		var fname = MSGID_FILE_NAME

		if item.MsgIdNameSpace != "" {
			namespace = item.MsgIdNameSpace
		}

		if item.MsgIdFileName != "" {
			fname = item.MsgIdFileName
		}

		//设置文件
		fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_CREATE)
		if err != nil {
			return err
		}

		//设置内容
		src := ""
		if isNewFile {
			//赋值定义
			src += "package " + namespace + "\n"
			src += "\n"
			src += "const ("
		} else {
			//读文件
			data, err := ioutil.ReadAll(fd)
			if err != nil {
				return err
			}
			content := string(data)
			i := strings.LastIndex(content, ")")
			if i > 0 {
				src += content[:i]
			} else {
				src += content
			}
		}
		fd.Close()
		src += item.FormatMsgId()
		src += ")\n"
		fd1, err := cleanFile(fname)
		if err != nil {
			return err
		}
		fd1.WriteString(src)
		fd1.Close()
	}
	return nil
}

// model*.go
func (gen *GoGenerator1) GenModel() error {
	//赋值
	items := SortWithId(gen.Items)
	for _, item := range items {
		var namespace = MODEL_PACK_NAME
		var fname = MODEL_FILE_NAME

		//设置文件
		if item.ModelNameSpace != "" {
			namespace = item.ModelNameSpace
		}
		if item.ModelFileName != "" {
			fname = item.ModelFileName
		}
		fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
		if err != nil {
			return err
		}
		src := ""
		if isNewFile {
			src = "package " + namespace + "\n"
		}
		src += item.FormatModel()
		fd.WriteString(src)
		fd.Close()
	}
	return nil
}

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
