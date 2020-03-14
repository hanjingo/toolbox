package http

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hanjingo/util"

	mm "github.com/hanjingo/component/memory_manager"
	fm "github.com/hanjingo/toolbox/file_master"
)

type Downloader1 struct {
	mu           *sync.Mutex
	filePathName string
	cache        chan *fm.FileSlice
	combineMap   map[int]*fm.FileSlice
	currIndex    int
	bDownload    bool
	memory       mm.MemoryManager1 //内存管理器
	finish       context.CancelFunc
	endTask      func(filePathName string)
}

func NewDownloader1(filePathName string, endTask func(filePathName string), memory mm.MemoryManager1) *Downloader1 {
	return &Downloader1{
		mu:           new(sync.Mutex),
		filePathName: filePathName,
		cache:        make(chan *fm.FileSlice, 1),
		combineMap:   make(map[int]*fm.FileSlice),
		currIndex:    1,
		bDownload:    false,
		memory:       memory,
		endTask:      endTask,
	}
}

func (down *Downloader1) Download() error {
	if down.bDownload {
		fmt.Printf("文件:%v已经被下载过了\n", down.filePathName)
		return errors.New("已经下载了")
	}
	down.bDownload = true

	var ctx context.Context
	ctx, down.finish = context.WithCancel(context.Background())
	down.combineAndWrite(ctx)
	return nil
}

func (down *Downloader1) Write(slice *fm.FileSlice) {
	go func() {
		select {
		case down.cache <- slice:
			return
		}
	}()
}

func (down *Downloader1) Finish(fileName string) {
	down.finish()
}

func (down *Downloader1) combineAndWrite(ctx context.Context) {
	go down.combine(ctx)
	go func() {
		for {
			select {
			case slice, ok := <-down.cache:
				if slice == nil && !ok {
					return
				}
				for !down.memory.Malloc(int64(slice.Buf.Len())) {
					time.Sleep(time.Duration(100 * time.Millisecond))
				}
				down.mu.Lock()
				down.combineMap[slice.Index] = slice
				down.mu.Unlock()
			}
		}
	}()
}

func (down *Downloader1) combine(ctx context.Context) error {
	defer func() {
		for len(down.cache) > 0 {
			slice := <-down.cache
			down.memory.Free(int64(slice.Buf.Len()))
		}
		close(down.cache)
		down.endTask(down.filePathName)
	}()
	f, err := os.OpenFile(down.filePathName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("打开文件失败,错误:%v\n", err)
		return err
	}
	defer f.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		if _, ok := down.combineMap[down.currIndex]; !ok {
			time.Sleep(time.Duration(100 * time.Millisecond))
			continue
		}
		slice := down.combineMap[down.currIndex]
		n := slice.Buf.Len()
		if _, err := f.Write(slice.GetContent()); err != nil {
			fmt.Printf("写入文件失败,错误:%v\n", err)
			return err
		}
		down.mu.Lock()
		delete(down.combineMap, down.currIndex)
		down.memory.Free(int64(n))
		fmt.Printf("内存:剩下%vMB\n", down.memory.Left()/util.MB)
		down.currIndex++
		down.mu.Unlock()
	}
	return nil
}
