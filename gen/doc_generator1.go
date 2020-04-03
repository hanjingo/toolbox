/*
doc.txt
"Name":"GenZhuReq",
"Desc":"跟注请求",
 "Fields":[
    ["凭证", "Token", "string"]
]

GenZhuReq
id: 12345
说明:
	跟住请求
参数		类型		说明
Token		string		凭证
示例:"{...}"

"Name":"ROOM_ERR",
"Desc":"房间通用错误",
"Id":1000

ROOM_ERR	id:1000    说明:房间错误
*/

package gen

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type DocGenerator1 struct {
	FileMap      map[string]string //文件列表
	NameSpaceMap map[string]string //命名空间
	Items        map[string]*Class //key:msg name  value:item集合
}

func NewDocGenerator1(cfgs []*ClassConfig, fileMap map[string]string,
	namespaceMap map[string]string) *DocGenerator1 {
	back := &DocGenerator1{
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

func (gen *DocGenerator1) Type() string {
	return LANG_DOC_V1
}

func (gen *DocGenerator1) Gen() error {
	if err := gen.GenModel(); err != nil {
		return err
	}
	if err := gen.GenErr(); err != nil {
		return err
	}
	return nil
}

func (gen *DocGenerator1) formatType(args ...string) string {
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
		return gen.formatType(args[1]) + "[]"
	case MAP:
		return "map[" + gen.formatType(args[1]) + "]" + gen.formatType(args[2])
	case POINT:
		return gen.formatType(args[1]) + "*"
	default:
		return args[0]
	}
}

//格式化示例
func (gen *DocGenerator1) formatExample(ci *Class) string {
	back := ""
	return back
}

// model*.doc
func (gen *DocGenerator1) formatModel(ci *Class) string {
	back := ""
	back += "\n"
	back += ci.Name + "\n"
	back += "消息id:" + strconv.Itoa(ci.Id) + "\n"
	back += formatLine(defSize, "参数", "类型", "说明") + "\n"
	for _, field := range ci.Fields {
		t := gen.formatType(field.Type, field.Type1, field.Type2)
		back += formatLine(defSize, field.Name, t, field.Desc) + "\n"
	}
	back += "说明:" + " " + ci.Desc + "\n"
	back += "示例:" + gen.formatExample(ci) + "\n"
	return back
}
func (gen *DocGenerator1) GenModel() error {
	if !isPrintDocModel(gen.FileMap) {
		return nil
	}
	//设置文件
	content := ""
	fname := gen.FileMap[KEY_DOC_MODEL]
	if fname == "" {
		return errors.New("生成doc版model时,无法读取空路径")
	}
	fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	defer fd.Close()
	if err != nil {
		return err
	}
	if !isNewFile {
		//读文件
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		content += string(data)
	} else {
		content += "消息模型:" + "\n"
	}
	for _, item := range gen.Items {
		content += gen.formatModel(item)
	}
	fd.WriteString(content)
	return nil
}

/*err*/
func (gen *DocGenerator1) formatErr(ci *Class) string {
	return formatLine(35, ci.Name, strconv.Itoa(ci.Id), ci.Desc) + "\n"
}
func (gen *DocGenerator1) GenErr() error {
	if !isPrintDocErr(gen.FileMap) {
		return nil
	}
	temp := make(map[string]*Class)
	for _, item := range gen.Items {
		temp[item.Name] = item
	}
	items := SortWithId(temp)

	//生成
	content := ""
	fname := gen.FileMap[KEY_DOC_ERR]
	if fname == "" {
		return errors.New("生成doc版errid时,无法读取空路径")
	}
	fd, isNewFile, err := mustOpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	defer fd.Close()
	if err != nil {
		return err
	}
	if !isNewFile {
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		content += string(data)
	}
	content += "\n定义错误码:" + "\n"
	//生成
	for _, item := range items {
		content += gen.formatErr(item)
	}
	fd.WriteString(content)
	return nil
}
