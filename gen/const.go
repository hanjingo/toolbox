package gen

type GenerI interface {
	Type() string
	Gen() error
}

const (
	LANG_GO     string = "GO"
	LANG_CSHARP string = "C#"
	LANG_JS     string = "JS"
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
	ACTION_PRINT_ID           string = "ID"
	ACTION_PRINT_MODEL        string = "MODEL"
	ACTION_PRINT_ERR          string = "ERR"
	ACTION_PRINT_ID_AND_MODEL string = "ID_AND_MODEL"
)

const (
	KEY_ID    string = "ID"
	KEY_MODEL string = "MODEL"
	KEY_ERR   string = "ERR"
)

var (
	ID_IDX  int = 0
	ERR_IDX int = 0
)
