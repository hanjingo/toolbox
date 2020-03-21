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
	Id           int               `json:"Id"`
	Desc         string            `json:"Desc"`
	Name         string            `json:"Name"`
	Fields       []*Field          `json:"Fields"`
	FileMap      map[string]string `json:"FileMap"`
	NameSpaceMap map[string]string `json:"NameSpaceMap"`
}

func NewClass(conf *ClassConfig, file map[string]string, namespace map[string]string) *Class {
	back := &Class{
		Fields:       []*Field{},
		FileMap:      make(map[string]string),
		NameSpaceMap: make(map[string]string),
	}
	back.Id = conf.Id
	back.Desc = conf.Desc
	back.Name = conf.Name
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
	for k, v := range file {
		back.FileMap[k] = v
	}
	for k, v := range namespace {
		back.NameSpaceMap[k] = v
	}
	return back
}
