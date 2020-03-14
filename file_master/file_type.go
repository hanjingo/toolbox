package file_master

import (
	"sync"
)

var contentTypeOnce = new(sync.Once)
var contentType map[string]string

func GetContentTypeMap() map[string]string {
	contentTypeOnce.Do(func() {
		contentType = make(map[string]string)
		doContentInit()
	})
	return contentType
}

func doContentInit() {
	contentType[".mp4"] = "video/mpeg4"
	contentType[".rmvb"] = "application/vnd.rn-realmedia-vbr"
}
