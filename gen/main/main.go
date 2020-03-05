package main

import (
	"fmt"
	"path/filepath"
	"strings"

	gg "github.com/hanjingo/toolbox/gen"
)

// for win:   go build -o gen.exe main.go
// for linux: go build -o gen main.go
func main() {
	for {
		var cmd string
		fmt.Println("请输入要生成的语言类型(go,js,c#...),按q退出>>")
		fmt.Scanln(&cmd)
		if strings.ToUpper(cmd) == "Q" {
			return
		}
		doGen(cmd)
	}
}

func doGen(lang string) {
	fmt.Println("请输入json文件绝对路径(例:c:\\007.json),默认读取当前路径的main.json>>")
	var addr string
	fmt.Scanln(&addr)
	if addr == "" {
		addr = filepath.Join(gg.GetCurrPath(), "main.json")
	}
	conf := gg.GetCodeGenConfig()
	err := gg.LoadJsonConfig(addr, conf)
	if err != nil {
		fmt.Println("加载json文件失败,错误:", err)
		return
	}
	gg.PATH = conf.Path
	gg.SetEnv()
	switch strings.ToUpper(lang) {
	case "GOLANG", "GO":
		gen := gg.NewGoGenerator1(conf)
		fmt.Println("生成消息id 结果:", check(gen.GenMsgid))
		fmt.Println("生成数据结构 结果:", check(gen.GenModel))
		fmt.Println("生成结束")
	case "JAVASCRIPT", "JS":
		return
	default:
		fmt.Println("不支持的语言类型")
		return
	}

}

func check(f func() error) string {
	err := f()
	if err != nil {
		return err.Error()
	}
	return "成功"
}
