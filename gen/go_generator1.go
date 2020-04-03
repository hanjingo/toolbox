package gen

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type GoGenerator1 struct {
	startIdx     int               //起始id
	startErrIdx  int               //起始错误id
	FileMap      map[string]string //文件列表
	NameSpaceMap map[string]string //命名空间
	Items        map[string]*Class //key:msg name  value:item集合
}

func NewGoGenerator1(cfgs []*ClassConfig, fileMap map[string]string,
	namespaceMap map[string]string) *GoGenerator1 {
	back := &GoGenerator1{
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

func (gen *GoGenerator1) Type() string {
	return LANG_GO_V1
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
		return errors.New("生成go版msgid时,无法读取空路径")
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
		//赋值定义
		content += "package " + namespace + "\n"
		content += "\n"
		content += "const ("
	} else {
		//读文件
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		temp := string(data)
		i := strings.LastIndex(temp, ")")
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
	content += "\n)"
	fd.WriteString(content)
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
	if !isPrintModel(gen.FileMap) {
		return nil
	}
	//设置文件
	content := ""
	fname := gen.FileMap[KEY_MODEL]
	if fname == "" {
		return errors.New("生成go版model时,无法读取空路径")
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
		content += "package " + namespace + "\n"
	} else {
		//读文件
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		content += string(data)
	}
	//赋值
	for _, item := range gen.Items {
		content += gen.formatModel(item)
	}
	fd.WriteString(content)
	return nil
}

/*errid.go*/
func (gen *GoGenerator1) formatErr(ci *Class) string {
	back := ""
	back += "\n"
	back += "	" + ci.Name
	back += " uint32 = " + strconv.Itoa(ci.Id) + "//" + ci.Desc
	return back
}
func (gen *GoGenerator1) GenErr() error {
	if !isPrintErr(gen.FileMap) {
		return nil
	}
	temp := make(map[string]*Class)
	for _, item := range gen.Items {
		if !isPrintErr(gen.FileMap) {
			continue
		}
		temp[item.Name] = item
	}
	items := SortWithId(temp)

	//开始填充数据
	content := ""
	fname := gen.FileMap[KEY_ERR]
	if fname == "" {
		return errors.New("生成go版errid时,无法读取空路径")
	}
	fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	defer fd.Close()
	if err != nil {
		return err
	}
	if isNewFile {
		//命名空间 or 包
		namespace := "err"
		if value, ok := gen.NameSpaceMap[KEY_ERR]; ok {
			namespace = value
		}
		content = "package " + namespace + "\n"
		content += "\n"
		content += "const ("
	} else {
		//读文件
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		temp := string(data)
		i := strings.LastIndex(temp, ")")
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
	content += "\n)"
	fd.WriteString(content)
	return nil
}
