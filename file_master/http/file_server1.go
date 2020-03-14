package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	mm "github.com/hanjingo/component/memory_manager"
	"github.com/hanjingo/network"
	fm "github.com/hanjingo/toolbox/file_master"
)

type FileServer1 struct {
	m      map[string]*Downloader1 //key: 文件名 value: 下载器
	ip     string
	port   int
	api    string
	memory *mm.MemoryManager1               //内存管理器
	path   string                           //文件存放路径
	f      func(args ...interface{}) string //文件路劲生成函数
}

func NewFileServer1(ip string, port int, api string, memory *mm.MemoryManager1, path string, f func(args ...interface{}) string) *FileServer1 {
	return &FileServer1{
		m:      make(map[string]*Downloader1),
		ip:     ip,
		port:   port,
		api:    api,
		memory: memory,
		path:   path,
		f:      f,
	}
}

func (s *FileServer1) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	http.HandleFunc(s.api, s.onApi)

	addr := fmt.Sprintf("%v:%v", s.ip, s.port)
	http.ListenAndServe(addr, nil)
}

func (s *FileServer1) onApi(w http.ResponseWriter, r *http.Request) {
	var filePathName string
	var index int
	var md5 string

	var data []byte
	var err error

	if err := r.ParseMultipartForm(s.memory.Left()); err != nil {
		fmt.Println("解析http失败,错误:", err)
		network.SendBack(w, "解析http失败", 405)
	}

	//文件名
	name := network.GetMultPartFormValue(r.MultipartForm, URL_FILE_NAME)
	filePathName = s.f(name)
	if filePathName == "" {
		fmt.Printf("解析http文件名失败\n")
		network.SendBack(w, "解析http文件名失败", 405)
		return
	}

	//索引
	if index, err = strconv.Atoi(network.GetMultPartFormValue(r.MultipartForm, URL_INDEX)); err != nil {
		fmt.Printf("解析http索引失败:%v\n", err)
		network.SendBack(w, "解析http索引失败", 405)
		return
	}
	slice := fm.NewFileSlice(filePathName, index)

	//md5
	md5 = network.GetMultPartFormValue(r.MultipartForm, URL_MD5)
	slice.MD5 = md5

	//文件
	data = network.GetMultPartFormData(r.MultipartForm, URL_FILE_DATA)
	if data == nil {
		fmt.Printf("解析http文件内容失败:%v\n", err)
		network.SendBack(w, "解析http文件内容失败", 405)
		return
	}

	if _, ok := s.m[filePathName]; !ok {
		if err := s.addTask(filePathName); err != nil {
			fmt.Printf("添加下载任务失败,错误:%v\n", err)
			return
		}
	}
	downloader := s.m[filePathName]

	slice.Buf.Write(data)
	if !slice.CheckMD5() {
		fmt.Printf("校验md5失败\n")
		network.SendBack(w, "校验md5失败", 405)
		return
	}

	downloader.Write(slice)
}

func (s *FileServer1) addTask(filePathName string) error {
	if filePathName == "" {
		return errors.New("文件绝对路径名不为空")
	}
	var downloader *Downloader1
	if _, ok := s.m[filePathName]; ok {
		return nil
	}
	downloader = NewDownloader1(filePathName, s.endTask, *s.memory)
	downloader.Download()
	s.m[filePathName] = downloader
	return nil
}

func (s *FileServer1) endTask(fileName string) {
	delete(s.m, fileName)
}
