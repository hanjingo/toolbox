package gen

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type GoGenerator1 struct {
	startIdx    int               //起始id
	startErrIdx int               //起始错误id
	Items       map[string]*Class //key:msg name  value:item集合
	Conf        *Config           //配置
}

func NewGoGenerator1(conf *Config, fileMap map[string]string,
	namespaceMap map[string]string) *GoGenerator1 {
	back := &GoGenerator1{
		Items: make(map[string]*Class),
		Conf:  conf,
	}
	for _, e := range conf.Classes {
		class := NewClass(e, fileMap, namespaceMap)
		if isPrintId(class.FileMap) {
			if class.Id != 0 {
				back.startIdx = class.Id
			} else {
				back.startIdx++
			}
			class.Id = back.startIdx
		}
		if isPrintErr(class.FileMap) {
			if class.Id != 0 {
				back.startErrIdx = class.Id
			} else {
				back.startErrIdx++
			}
			class.Id = back.startErrIdx
		}
		back.Items[class.Name] = class
	}
	return back
}

func (gen *GoGenerator1) Type() string {
	return LANG_GO
}

func (gen *GoGenerator1) Gen() error {
	if err := gen.GenMsgid(); err != nil {
		return err
	}
	if err := gen.GenModel(); err != nil {
		return err
	}
	if err := gen.GenErr(); err != nil {
		return err
	}
	return nil
}

func (gen *GoGenerator1) formatType(args ...string) string {
	if args == nil || len(args) == 0 {
		return ""
	}
	switch strings.ToUpper(args[0]) {
	case UINT8:
		return "uint8"
	case UINT32:
		return "uint32"
	case UINT64:
		return "uint64"
	case INT:
		return "int"
	case INT64:
		return "int64"
	case FLOAT:
		return "float32"
	case DOUBLE:
		return "float64"
	case STRING:
		return "string"
	case BOOL:
		return "bool"
	case ARRAY:
		return "[]" + gen.formatType(args[1])
	case MAP:
		return "map[" + gen.formatType(args[1]) + "]" + gen.formatType(args[2])
	case POINT:
		return "*" + gen.formatType(args[1])
	default:
		return args[0]
	}
}

// msgid.go
func (gen *GoGenerator1) formatMsgid(ci *Class) string {
	var back = ""
	back += "\n"
	back += "    " + ci.Name + " uint32 = " + strconv.Itoa(ci.Id)
	if ci.Desc != "" {
		back += " " + "//" + ci.Desc
	}
	return back
}
func (gen *GoGenerator1) GenMsgid() error {
	//msgid
	temp := make(map[string]*Class)
	for _, item := range gen.Items {
		if !isPrintId(item.FileMap) {
			continue
		}
		temp[item.Name] = item
	}
	items := SortWithId(temp)
	for _, item := range items {
		namespace := ""
		if value, ok := item.NameSpaceMap[KEY_ID]; ok {
			namespace = value
		}

		//设置文件
		fname := item.FileMap[KEY_ID]
		fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
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
		src += gen.formatMsgid(item)
		src += "\n)"
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
func (gen *GoGenerator1) formatModel(ci *Class) string {
	var back = ""
	back += "\n"
	if ci.Desc != "" {
		back += "//" + ci.Desc + "\n"
	}
	back += "type " + ci.Name + " struct {\n"
	for _, field := range ci.Fields {
		back += "	" + field.Name + " "
		back += gen.formatType(field.Type, field.Type1, field.Type2)
		if field.Desc != "" {
			back += " " + "//" + field.Desc
		}
		back += "\n"
	}
	back += "}\n"
	return back
}
func (gen *GoGenerator1) GenModel() error {
	//赋值
	for _, item := range gen.Items {
		if !isPrintModel(item.FileMap) {
			continue
		}

		//设置文件
		namespace := ""
		if value, ok := item.NameSpaceMap[KEY_MODEL]; ok {
			namespace = value
		}
		fname := item.FileMap[KEY_MODEL]
		fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
		if err != nil {
			return err
		}
		src := ""
		if isNewFile {
			src += "package " + namespace + "\n"
		}
		src += gen.formatModel(item)
		fd.WriteString(src)
		fd.Close()
	}
	return nil
}

/*errid.go*/
func (gen *GoGenerator1) formatErr(ci *Class) string {
	back := ""
	back += "\n"
	back += "	" + ci.Name
	back += " = " + strconv.Itoa(ci.Id) + "//" + ci.Desc
	return back
}
func (gen *GoGenerator1) GenErr() error {
	temp := make(map[string]*Class)
	for _, item := range gen.Items {
		if !isPrintErr(item.FileMap) {
			continue
		}
		temp[item.Name] = item
	}
	items := SortWithId(temp)
	//生成
	for _, item := range items {
		//命名空间 or 包
		namespace := "err"
		if value, ok := item.NameSpaceMap[KEY_ERR]; ok {
			namespace = value
		}
		fname := item.FileMap[KEY_ERR]
		fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
		if err != nil {
			return err
		}
		src := ""
		if isNewFile {
			src = "package " + namespace + "\n"
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
		src += gen.formatErr(item)
		src += "\n)"
		fd1, err := cleanFile(fname)
		if err != nil {
			return err
		}
		fd1.WriteString(src)
		fd1.Close()
	}
	return nil
}
