package gen

import (
	"io/ioutil"
	"os"
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
	for _, e := range conf.Classes {
		class := NewClass(e, conf.PathMap, conf.NameSpaceMap)
		back.Items[class.Name] = class
	}
	return back
}

func (gen *CsGenerator1) Type() string {
	return LANG_GO
}

func (gen *CsGenerator1) Gen() error {
	if err := gen.GenMsgid(); err != nil {
		return err
	}
	if err := gen.GenModel(); err != nil {
		return err
	}
	return nil
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
		return gen.formatType(args[1]) + "*"
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
	start := 0
	temp := make(map[string]*Class)
	for _, item := range gen.Items {
		if !isPrintId(item.FileMap) {
			continue
		}
		if item.Id != 0 {
			start = item.Id
		}
		item.Id = start
		start++
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
			if namespace != "" {
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

/*errid.cs*/
func (gen *CsGenerator1) formatErr(ci *Class) string {
	back := ""
	back += "\n"
	back += "	" + ci.Name
	back += " = err(" + strconv.Itoa(ci.Id) + ", " + "\"" + ci.Desc + "\")"
	return back
}
func (gen *CsGenerator1) GenErr() error {
	start := 0
	temp := make(map[string]*Class)
	for _, item := range gen.Items {
		if !isPrintErr(item.FileMap) {
			continue
		}
		if item.Id != 0 {
			start = item.Id
		}
		item.Id = start
		start++
		temp[item.Name] = item
	}
	items := SortWithId(temp)
	//生成
	for _, item := range items {
		//命名空间 or 包
		namespace := "errid"
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
			//赋值定义
			src += "using System;\n"
			src += "using System.Collections;\n"
			src += "using System.Collections.Generic;\n"
			src += "\n"
			if namespace != "" {
				src += "namespace " + namespace + "{"
			}
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
