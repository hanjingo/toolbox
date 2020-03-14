package main

import (
	"path/filepath"
	"sync"

	"github.com/hanjingo/util"

	mm "github.com/hanjingo/component/memory_manager"
	fs "github.com/hanjingo/toolbox/file_master/http"
)

func genFilePathName(args ...interface{}) string {
	if args == nil || len(args) == 0 {
		return ""
	}
	name := args[0].(string)
	path := util.GetCurrPath()
	return filepath.Join(path, name)
}

// for win:   go build -o fs.exe main.go
// for linux: set GOARCH=386
// for linux: set GOOS=linux
// for linux: go build -o fs main.go
func main() {
	wg := new(sync.WaitGroup)
	memory := mm.NewMemoryManager1(100 * util.MB)
	s := fs.NewFileServer1("127.0.0.1", 10086, "/upload", memory, util.GetCurrPath(), genFilePathName)
	wg.Add(1)
	go s.Run(wg)
	wg.Wait()
}
