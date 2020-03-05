package gen

import (
	"fmt"
)

type Field struct {
	Desc  string `json:"Desc"`
	Name  string `json:"Name"`
	Type  string `json:"Type"`
	Type1 string `json:"Type1"`
	Type2 string `json:"Type2"`
}

type Class struct {
	NoPrintId      bool     `json:"NoPrintId"`
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
	back.NoPrintId = conf.NoPrintId
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
