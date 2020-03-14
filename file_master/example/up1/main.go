package main

import (
	"fmt"
	"time"

	mm "github.com/hanjingo/component/memory_manager"
	fs "github.com/hanjingo/toolbox/file_master/http"
	"github.com/hanjingo/util"
)

var lastIndex float64 = 0

func doUpdate(progress float64) {
	if progress-lastIndex > 0.1 {
		fmt.Print("*")
		lastIndex = progress
	}
}

// for win:   go build -o uploader.exe main.go
// for linux: set GOARCH=386
// for linux: set GOOS=linux
// for linux: go build -o uploader main.go
func main() {
	for {
		var path string
		var url string
		fmt.Println("请输入要上传的文件绝对路径>>")
		fmt.Scanln(&path)
		memory := mm.NewMemoryManager1(100 * util.MB)
		uploader := fs.NewUploader1(1*util.MB, memory, time.Duration(10*time.Second))
		fmt.Println("请输入要上传到的服务器url>>")
		fmt.Scanln(&url)
		uploader.Upload(path, url)
		fmt.Print("进度:[")
		for {
			progress := uploader.GetProgress()
			if progress < 1 {
				time.Sleep(time.Duration(100 * time.Millisecond))
				doUpdate(progress)
			} else {
				fmt.Print("]\n")
				goto doEnd
			}
		}
	doEnd:
		fmt.Println("上传完成！！！")
		fmt.Println("")
	}
}
