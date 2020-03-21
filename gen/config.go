package gen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type MainConfig struct {
	PathMap   map[string]string  `json:"PathMap"`
	configMap map[string]*Config //配置集合
}

func (mc *MainConfig) Load(filePath string) error {
	if err := LoadJsonConfig(filePath, mc); err != nil {
		return err
	}
	for name, str := range mc.PathMap {
		conf := NewConfig()
		if err := LoadJsonConfig(str, conf); err != nil {
			return err
		}
		conf.Init()
		mc.configMap[name] = conf
	}
	return nil
}

func NewMainConfig() *MainConfig {
	return &MainConfig{
		PathMap:   make(map[string]string),
		configMap: make(map[string]*Config),
	}
}

type Config struct {
	Lang         string            `json:"Lang"`      //要生成的语言类型 js,c#,go...
	PathMap      map[string]string `json:"Path"`      //要生成的文件存放的路径
	NameSpaceMap map[string]string `json:"NameSpace"` //命名空间配置
	Classes      []*ClassConfig    `json:"Class"`
}

func (cfg *Config) Init() {
	for k, v := range cfg.PathMap {
		switch strings.ToUpper(cfg.Lang) {
		case LANG_GO:
			if ext := filepath.Ext(v); ext == "" {
				cfg.PathMap[k] = v + ".go"
			}
		case LANG_CSHARP:
			if ext := filepath.Ext(v); ext == "" {
				cfg.PathMap[k] = v + ".cs"
			}
		case LANG_JS:
		}
	}
}

type ClassConfig struct {
	Id     int        `json:"Id"`
	Name   string     `json:"Name"`
	Desc   string     `json:"Desc"`
	Fields [][]string `json:"Fields"`
}

func NewConfig() *Config {
	back := &Config{
		PathMap:      make(map[string]string),
		NameSpaceMap: make(map[string]string),
		Classes:      []*ClassConfig{},
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
