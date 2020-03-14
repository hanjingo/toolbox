package file_master

import "sync"

type UploaderI interface {
	Upload(filePath, url string) error //上传
	GetProgress() float64              //进度
	Cancel() error                     //取消
}

type DownloaderI interface {
	Download() error        //下载
	Write(slice *FileSlice) //写入
	Finish(fileName string) //结束下载任务
}

type FileServerI interface {
	Run(wg *sync.WaitGroup) //文件服务器跑起来(阻塞)
}
