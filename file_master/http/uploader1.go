package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hanjingo/util"

	fm "github.com/hanjingo/component/file_manager"
	mm "github.com/hanjingo/component/memory_manager"
)

type Uploader1 struct {
	memory         mm.MemoryManagerI  //内存管理器
	maxSliceLimit  int64              //最大切片上传限制(B)
	uploadTimeOut  time.Duration      //上传超时 单位(s)
	cache          chan *fm.FileSlice //上传缓存
	bUpload        bool               //是否已经上传
	total_slice    float64            //总片数
	uploaded_slice float64            //已经上传的片数
}

func NewUploader1(limit int64, memory mm.MemoryManagerI, uploadDur time.Duration) *Uploader1 {
	return &Uploader1{
		memory:         memory,
		maxSliceLimit:  limit,
		uploadTimeOut:  uploadDur,
		cache:          make(chan *fm.FileSlice, 1),
		bUpload:        false,
		total_slice:    0.0,
		uploaded_slice: 0.0,
	}
}

func (up *Uploader1) Upload(filePath, url string) error {
	if up.bUpload {
		fmt.Printf("文件已经上传了\n")
		return errors.New("文件已经上传了")
	}
	up.bUpload = true
	if !up.check() {
		return errors.New("配置错误")
	}
	go up.readAndSlice(filePath)
	go func() {
		tm := time.NewTimer(up.uploadTimeOut)
		for {
			select {
			case <-tm.C:
				fmt.Printf("\n上传超时\n")
			case arg := <-up.cache:
				go up.doUpload(url, arg)
			}
		}
	}()

	return nil
}

func (up *Uploader1) GetProgress() float64 {
	if up.total_slice == 0.0 {
		return 0.0
	}
	return up.uploaded_slice / up.total_slice
}

func (up *Uploader1) Cancel() error {
	//todo
	return nil
}

func (up *Uploader1) check() bool {
	if up.maxSliceLimit <= 0 {
		return false
	}
	if up.cache == nil {
		return false
	}
	return true
}

func (up *Uploader1) doUpload(url string, slice *fm.FileSlice) error {
	defer func() {
		up.memory.Free(up.maxSliceLimit)
	}()
	fmt.Printf("上传切片 name:%v, index:%v\n", slice.Name, slice.Index)

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	writer.WriteField(URL_FILE_NAME, slice.Name)
	writer.WriteField(URL_INDEX, strconv.Itoa(slice.Index))
	writer.WriteField(URL_MD5, slice.MD5)

	formFile, err := writer.CreateFormFile(URL_FILE_DATA, slice.Name)
	if err != nil {
		fmt.Printf("创建表单文件失败,错误:%v\n", err)
		return err
	}
	if _, err := formFile.Write(slice.Buf.Bytes()); err != nil {
		fmt.Printf("复制切片信息失败")
		return err
	}
	contentType := writer.FormDataContentType()
	writer.Close()

	//发送
	rsp, err := http.Post(url, contentType, buf)
	if err != nil {
		fmt.Printf("上传切片失败,错误:%v\n", err)
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode >= 200 && rsp.StatusCode < 300 {
		up.uploaded_slice++
	} else {
		fmt.Println("上传失败,错误:", rsp.Status)
	}
	return nil
}

func (up *Uploader1) readAndSlice(filePath string) error {
	size, err := util.GetFileSize(filePath, util.BYTE)
	if err != nil {
		fmt.Println("获得文件大小失败,错误:%v", err)
		return err
	}

	name := util.GetFileFullName(filePath)
	slice_num := int((size / up.maxSliceLimit) + 1)
	up.total_slice = float64(slice_num)
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("打开文件失败,错误:", err)
		return err
	}
	defer f.Close()
	for i := 1; ; i++ {
		for !up.memory.Malloc(up.maxSliceLimit) {
			time.Sleep(time.Duration(100 * time.Millisecond))
		}
		slice := fm.NewFileSlice(name, i)
		var temp = make([]byte, up.maxSliceLimit)
		n, err := f.Read(temp)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("读取文件:%v 失败,错误:%v", filePath, err)
				return err
			}
		}
		slice.Buf.Write(temp[:n])
		slice.MD5 = slice.CalcMD5()
		if n > 0 {
			up.cache <- slice
		}
		if err == io.EOF {
			return nil
		}
	}
	return nil
}
