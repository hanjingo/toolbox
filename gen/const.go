package gen

type GenerI interface {
	Type() string
	Gen() error
}

//生成器版本
const (
	LANG_GO_V1     string = "GO_V1"
	LANG_CSHARP_V1 string = "C#_V1"
	LANG_JS_V1     string = "JS_V1"
	LANG_DOC_V1    string = "DOC_V1"
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
	POINT  string = "*"
)

const (
	KEY_ID        string = "ID"
	KEY_MODEL     string = "MODEL"
	KEY_ERR       string = "ERR"
	KEY_DOC_MODEL string = "DOC_MODEL"
	KEY_DOC_ERR   string = "DOC_ERR"
)

const defSize int = 24 //列宽
