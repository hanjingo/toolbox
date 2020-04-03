package gen

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type CsGenerator1 struct {
	FileMap      map[string]string //文件列表
	NameSpaceMap map[string]string //命名空间
	Items        map[string]*Class //key:msg value:item集合
}

func NewCsGenerator1(cfgs []*ClassConfig, fileMap map[string]string,
	namespaceMap map[string]string) *CsGenerator1 {
	back := &CsGenerator1{
		FileMap:      fileMap,
		NameSpaceMap: namespaceMap,
		Items:        make(map[string]*Class),
	}
	for _, e := range cfgs {
		class := NewClass(e)
		back.Items[e.Name] = class
	}
	return back
}

func (gen *CsGenerator1) Type() string {
	return LANG_CSHARP_V1
}

func (gen *CsGenerator1) Gen() error {
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
		return "int"
	case INT64:
		return "Int64"
	case FLOAT:
		return "Float32"
	case DOUBLE:
		return "Float64"
	case STRING:
		return "string"
	case BOOL:
		return "bool"
	case ARRAY:
		return gen.formatType(args[1]) + "[]"
	case MAP:
		return "Dictionary<" + gen.formatType(args[1]) + ", " + gen.formatType(args[2]) + ">"
	case POINT:
		return gen.formatType(args[1])
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
	if !isPrintId(gen.FileMap) {
		return nil
	}
	//msgid
	temp := make(map[string]*Class)
	for _, item := range gen.Items {
		if !isPrintId(gen.FileMap) {
			continue
		}
		temp[item.Name] = item
	}
	items := SortWithId(temp)
	//读文件先
	fname := gen.FileMap[KEY_ID]
	if fname == "" {
		return errors.New("生成c#版msgid时,无法读取空路径")
	}
	fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	defer fd.Close()
	if err != nil {
		return err
	}

	content := ""
	if isNewFile {
		namespace := ""
		if value, ok := gen.NameSpaceMap[KEY_ID]; ok {
			namespace = value
		}
		content += "using System;\n"
		content += "\n"
		if namespace != "" {
			content += "public enum " + namespace + ":UInt32{"
		}
	} else {
		//读文件
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		temp := string(data)
		i := strings.LastIndex(temp, "}")
		if i > 0 {
			content += temp[:i]
		} else {
			content += temp
		}
	}
	//设置内容
	for _, item := range items {
		content += gen.formatMsgid(item)
	}
	content += "\n}"
	fd.WriteString(content)
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
	if !isPrintModel(gen.FileMap) {
		return nil
	}
	//设置文件
	content := ""
	fname := gen.FileMap[KEY_MODEL]
	if fname == "" {
		return errors.New("生成c#版model时,无法读取空路径")
	}
	fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	defer fd.Close()
	if err != nil {
		return err
	}
	if isNewFile {
		namespace := ""
		if value, ok := gen.NameSpaceMap[KEY_MODEL]; ok {
			namespace = value
		}
		//赋值定义
		content += "using System;\n"
		content += "\n"
		if namespace != "" {
			content += "namespace " + namespace + "{\n"
		}
	} else {
		//读文件
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		temp := string(data)
		i := strings.LastIndex(temp, "}")
		if i > 0 {
			content += temp[:i]
		} else {
			content += temp
		}
	}
	//赋值
	items := SortWithId(gen.Items)
	for _, item := range items {
		content += gen.formatModel(item)
	}
	content += "\n}"
	fd.WriteString(content)
	return nil
}

/*errid.cs*/
func (gen *CsGenerator1) formatErr(ci *Class) string {
	back := ""
	back += "\n"
	back += "	" + ci.Name
	back += " = " + strconv.Itoa(ci.Id) + ", //" + ci.Desc
	return back
}
func (gen *CsGenerator1) GenErr() error {
	if !isPrintErr(gen.FileMap) {
		return nil
	}
	temp := make(map[string]*Class)
	for _, item := range gen.Items {
		temp[item.Name] = item
	}
	items := SortWithId(temp)

	//开始填充数据
	content := ""
	fname := gen.FileMap[KEY_ERR]
	if fname == "" {
		return errors.New("生成c#版errid时,无法读取空路径")
	}
	fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	defer fd.Close()
	if err != nil {
		return err
	}
	if isNewFile {
		//赋值定义
		namespace := "Err"
		if value, ok := gen.NameSpaceMap[KEY_ERR]; ok {
			namespace = value
		}
		content += "using System;\n"
		content += "\n"
		if namespace != "" {
			content += "public enum " + namespace + ":UInt32{"
		}
	} else {
		//读文件
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		temp := string(data)
		i := strings.LastIndex(temp, "}")
		if i > 0 {
			content += temp[:i]
		} else {
			content += temp
		}
	}

	//生成
	for _, item := range items {
		content += gen.formatErr(item)
	}
	content += "\n}"
	fd.WriteString(content)
	return nil
}
