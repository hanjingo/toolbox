package gen

import (
	"path/filepath"
)

const (
	UINT8  string = "UINT8"
	UINT32 string = "UINT32"
	UINT64 string = "UINT64"
	INT    string = "INT"
	INT64  string = "INT64"
	FLOAT  string = "FLOAT"
	DOUBLE string = "DOUBLE"
	STRING string = "STRING"
	BOOL   string = "BOOL"
	ARRAY  string = "ARRAY"
	MAP    string = "MAP"
)

var (
	PATH = GetCurrPath()

	MSGID_PACK_NAME = "msgid"
	MSGID_FILE_NAME = ""

	MODEL_PACK_NAME = "model"
	MODEL_FILE_NAME = ""
)

func SetEnv() {
	MSGID_FILE_NAME = filepath.Join(PATH, MSGID_PACK_NAME)
	MODEL_FILE_NAME = filepath.Join(PATH, MODEL_PACK_NAME)
}
