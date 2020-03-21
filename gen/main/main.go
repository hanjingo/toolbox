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
		var addr string
		fmt.Println("请输入json文件绝对路径(例:c:\\007.json),按q退出,默认读取当前路径的main.json>>")
		fmt.Scanln(&addr)
		if strings.ToUpper(addr) == "Q" {
			return
		}
		doGen(addr)
	}
}

func doGen(addr string) {
	if addr == "" {
		addr = filepath.Join(gg.GetCurrPath(), "main.json")
	}
	app := gg.GetApp()
	if err := app.Load(addr); err != nil {
		fmt.Println("加载文件:", addr, "失败")
		return
	}
	if err := app.Gen(); err != nil {
		fmt.Println("生成失败,错误:", err)
		return
	}
	fmt.Println("生成成功!!!")
}
