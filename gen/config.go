package gen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
)

type Config struct {
	GameName   string         `json:"GameName"`
	Path       string         `json:"Path"`
	StartIndex int            `json:"StartIndex"`
	Classes    []*ClassConfig `json:"Class"`
}

type ClassConfig struct {
	Id             int        `json:"Id"`
	Name           string     `json:"Name"`
	Desc           string     `json:"Desc"`
	MsgIdFileName  string     `json:"MsgIdFileName"`
	MsgIdNameSpace string     `json:"MsgIdNameSpace"` //msgid 命名空间
	ModelFileName  string     `json:"ModelFileName"`
	ModelNameSpace string     `json:"ModelNameSpace"` //model 命名空间
	Fields         [][]string `json:"Fields"`
}

var cgc *Config
var cgcOnce = new(sync.Once)

func GetCodeGenConfig() *Config {
	cgcOnce.Do(func() {
		cgc = &Config{Classes: []*ClassConfig{}}
	})
	return cgc
}

/*加载json格式的配置文件*/
func LoadJsonConfig(file_path string, conf interface{}) error {
	if file_path == "" {
		fmt.Println("配置文件地址不能为空")
		return errors.New("配置文件地址不能为空")
	}
	data, err := ioutil.ReadFile(file_path)
	if err != nil {
		fmt.Printf("加载配置文件:%v失败,错误:%v\n", file_path, err)
		return err
	}
	err = json.Unmarshal(data, conf)
	if err != nil {
		fmt.Printf("解析配置文件:%v失败,错误:%v\n", file_path, err)
		return err
	}
	fmt.Printf("配置文件:%v加载成功\n", file_path)
	return nil
}
