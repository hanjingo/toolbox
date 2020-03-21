package gen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type MainConfig struct {
	FileMap   map[string]*FileConfig `json:"File"`
	configMap map[string]*Config     //配置集合
}
type FileConfig struct {
	Lang         string            `json:"Lang"`       //要生成的语言类型 js,c#,go...
	Path         string            `json:"Path"`       //文件路劲
	PathMap      map[string]string `json:"CreatePath"` //要生成的文件存放的路径
	NameSpaceMap map[string]string `json:"NameSpace"`  //命名空间配置
}

func (mc *MainConfig) Load(filePath string) error {
	if err := LoadJsonConfig(filePath, mc); err != nil {
		return err
	}
	for name, f := range mc.FileMap {
		conf := NewConfig()
		if err := LoadJsonConfig(f.Path, conf); err != nil {
			return err
		}
		mc.configMap[name] = conf
	}
	//如果要创建的文件已经存在，删掉它
	for _, cfg := range mc.FileMap {
		for _, addr := range cfg.PathMap {
			fmt.Println("文件:", addr, " 被删除")
			os.Remove(addr)
		}
	}
	return nil
}

func NewMainConfig() *MainConfig {
	return &MainConfig{
		FileMap:   make(map[string]*FileConfig),
		configMap: make(map[string]*Config),
	}
}

type Config struct {
	Classes []*ClassConfig `json:"Class"`
}

type ClassConfig struct {
	Id     int        `json:"Id"`
	Name   string     `json:"Name"`
	Desc   string     `json:"Desc"`
	Fields [][]string `json:"Fields"`
}

func NewConfig() *Config {
	back := &Config{
		Classes: []*ClassConfig{},
	}
	return back
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
