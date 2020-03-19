package gen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CsGenerator1 struct {
	Items map[string]*Class //key:msg value:item集合
}

func NewCsGenerator1(conf *Config) *CsGenerator1 {
	back := &CsGenerator1{
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

func (gen *CsGenerator1) formatType(args ...string) string {
	if args == nil || len(args) == 0 {
		return ""
	}
	switch strings.ToUpper(args[0]) {
	case UINT8:
		return "UInt8"
	case UINT32:
		return "UInt32"
	case UINT64:
		return "UInt64"
	case INT:
		return "Int"
	case INT64:
		return "Int64"
	case FLOAT:
		return "Float32"
	case DOUBLE:
		return "Float64"
	case STRING:
		return "String"
	case BOOL:
		return "Bool"
	case ARRAY:
		return gen.formatType(args[1]) + "[]"
	case MAP:
		return "Dictionary<" + gen.formatType(args[1]) + ", " + gen.formatType(args[2]) + ">"
	case POINT:
		return gen.formatType(args[1]) + "*";
	default:
		return args[0]
	}
}

// 生成msgid.cs文件
func (gen *CsGenerator1) formatMsgid(ci *Class) string {
	var back = ""
	back += "\n"
	back += "    " + ci.Name + " = " + strconv.Itoa(ci.Id) + ","
	if ci.Desc != "" {
		back += " " + "//" + ci.Desc
	}
	return back
}
func (gen *CsGenerator1) GenMsgid() error {
	//msgid
	items := SortWithId(gen.Items)
	for _, item := range items {
		default_file_name := "msgid.cs"
		if item.NoPrintId {
			continue
		}
		if item.MsgIdFileName != "" {
			default_file_name = item.MsgIdFileName + ".cs"
		}
		namespace := "MsgId"
		fname := filepath.Join(MSGID_PATH_NAME, default_file_name)

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
			src += "using System;\n"
			src += "using System.Collections;\n"
			src += "using System.Collections.Generic;\n"
			src += "\n"
			if(namespace != "") {
				src += "public enum " + namespace + ":UInt32{"
			}
		} else {
			//读文件
			data, err := ioutil.ReadAll(fd)
			if err != nil {
				return err
			}
			content := string(data)
			i := strings.LastIndex(content, "}")
			if i > 0 {
				src += content[:i]
			} else {
				src += content
			}
		}
		fd.Close()
		src += gen.formatMsgid(item)
		src += "\n}"
		fd1, err := cleanFile(fname)
		if err != nil {
			return err
		}
		fd1.WriteString(src)
		fd1.Close()
	}
	return nil
}

// model*.cs
func (gen *CsGenerator1) formatModel(ci *Class) string {
	var back = ""
	back += "\n"
	if ci.Desc != "" {
		back += "//" + ci.Desc + "\n"
	}
	back += "public class " + ci.Name + " {\n"
	for _, field := range ci.Fields {
		back += "	public "
		back += gen.formatType(field.Type, field.Type1, field.Type2)
		back += " " + field.Name + ";"
		if field.Desc != "" {
			back += " " + "//" + field.Desc
		}
		back += "\n"
	}
	back += "}"
	return back
}
func (gen *CsGenerator1) GenModel() error {
	//赋值
	items := SortWithId(gen.Items)
	for _, item := range items {
		default_file_name := "model.cs"
		if item.ModelFileName != "" {
			default_file_name = item.ModelFileName + ".cs"
		}
		var namespace = "Model"
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
			//赋值定义
			src += "using System;\n"
			src += "using System.Collections;\n"
			src += "using System.Collections.Generic;\n"
			src += "\n"
			if namespace != "" {
				src = "namespace " + namespace + "{\n"
			}
		} else {
			//读文件
			data, err := ioutil.ReadAll(fd)
			if err != nil {
				return err
			}
			content := string(data)
			i := strings.LastIndex(content, "}")
			if i > 0 {
				src += content[:i]
			} else {
				src += content
			}
		}
		fd.Close()
		src += gen.formatModel(item)
		src += "\n}"
		fd1, err := cleanFile(fname)
		if err != nil {
			return err
		}
		fd1.WriteString(src)
		fd1.Close()
	}
	return nil
}