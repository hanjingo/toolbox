package gen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type GoGenerator1 struct {
	Items map[string]*Class //key:msg name  value:item集合
}

func NewGoGenerator1(conf *Config) *GoGenerator1 {
	back := &GoGenerator1{
		Items: make(map[string]*Class),
	}
	start := conf.StartIndex
	for _, e := range conf.Classes {
		if !e.NoPrintId {
			if e.Id != 0 {
				start = e.Id
			} else {
				start++
			}
			e.Id = start
		}
		class := NewClass(e)
		back.Items[class.Name] = class
	}
	return back
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
	items := SortWithId(gen.Items)
	for _, item := range items {
		default_file_name := "msgid.go"
		if item.NoPrintId {
			continue
		}
		if item.MsgIdFileName != "" {
			default_file_name = item.MsgIdFileName + ".go"
		}
		var namespace = MSGID_PACK_NAME
		var fname = filepath.Join(MSGID_PATH_NAME, default_file_name)

		if item.MsgIdNameSpace != "" {
			namespace = item.MsgIdNameSpace
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
	items := SortWithId(gen.Items)
	for _, item := range items {
		default_file_name := "model.go"
		if item.ModelFileName != "" {
			default_file_name = item.ModelFileName + ".go"
		}
		var namespace = MODEL_PACK_NAME
		var fname = filepath.Join(MODEL_PATH_NAME, default_file_name)

		//设置文件
		if item.ModelNameSpace != "" {
			namespace = item.ModelNameSpace
		}
		fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
		if err != nil {
			return err
		}
		src := ""
		if isNewFile {
			src = "package " + namespace + "\n"
		}
		src += gen.formatModel(item)
		fd.WriteString(src)
		fd.Close()
	}
	return nil
}

/*handle.go*/
