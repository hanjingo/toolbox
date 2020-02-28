package gen

import (
	"fmt"
	"strconv"
	"strings"
)

type Field struct {
	Desc  string `json:"Desc"`
	Name  string `json:"Name"`
	Type  string `json:"Type"`
	Type1 string `json:"Type1"`
	Type2 string `json:"Type2"`
}

type Class struct {
	Id             int      `json:"Id"`
	Desc           string   `json:"Desc"`
	Name           string   `json:"Name"`
	MsgIdFileName  string   `json:"MsgIdFileName"`
	MsgIdNameSpace string   `json:"MsgIdNameSpace"`
	ModelFileName  string   `json:"ModelFileName"`
	ModelNameSpace string   `json:"ModelNameSpace"`
	Fields         []*Field `json:"Fields"`
}

func NewClass(conf *ClassConfig) *Class {
	back := &Class{Fields: []*Field{}}
	back.Id = conf.Id
	back.Desc = conf.Desc
	back.Name = conf.Name
	back.MsgIdFileName = conf.MsgIdFileName
	back.MsgIdNameSpace = conf.MsgIdNameSpace
	back.ModelFileName = conf.ModelFileName
	back.ModelNameSpace = conf.ModelNameSpace
	for _, value := range conf.Fields {
		if value == nil || len(value) < 3 {
			fmt.Println("错误的类定义:", conf.Name)
			continue
		}
		field := &Field{
			Desc: value[0],
			Name: value[1],
		}
		switch len(value) {
		case 3:
			field.Type = value[2]
		case 4:
			field.Type = value[2]
			field.Type1 = value[3]
		case 5:
			field.Type = value[2]
			field.Type1 = value[3]
			field.Type2 = value[4]
		}
		back.Fields = append(back.Fields, field)
	}
	return back
}

//写 msgid
func (ci *Class) FormatMsgId() string {
	var back = ""
	back += "\n"
	back += "    " + ci.Name + "Req " + "uint32 = " + strconv.Itoa(ci.Id)
	if ci.Desc != "" {
		back += "    " + "//" + ci.Desc
	}
	back += "\n"
	return back
}

func formatType(args ...string) string {
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
		return "[]" + formatType(args[0])
	case MAP:
		return "map[" + formatType(args[1]) + "]" + formatType(args[2])
	default:
		return ""
	}
}

//写 model
func (ci *Class) FormatModel() string {
	var back = ""
	back += "\n"
	if ci.Desc != "" {
		back += "//" + ci.Desc + "\n"
	}
	back += "type " + ci.Name + " struct {\n"
	for _, field := range ci.Fields {
		back += field.Name + " "
		back += formatType(field.Type, field.Type1, field.Type2)
	}
	back += "}\n"
	return back
}
