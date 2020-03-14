package file_master

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
)

type FileSlice struct {
	Buf         *bytes.Buffer
	Index       int
	ContentType string
	Name        string
	MD5         string
}

func NewFileSlice(name string, index int) *FileSlice {
	return &FileSlice{
		Buf:   new(bytes.Buffer),
		Index: index,
		Name:  name,
	}
}

//拿到文件内容
func (fs *FileSlice) GetContent() []byte {
	return fs.Buf.Bytes()
}

//计算md5 base64
func (fs *FileSlice) CalcMD5() string {
	hash := md5.New()
	back := base64.StdEncoding.EncodeToString(hash.Sum(fs.Buf.Bytes()))
	return back
}

//校验md5
func (fs *FileSlice) CheckMD5() bool {
	return fs.CalcMD5() == fs.MD5
}
