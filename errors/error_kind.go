package errors

type Kind int

const (
	KindConfig     Kind = iota
	KindParse      Kind = iota
	KindFileSystem Kind = iota
	KindDB         Kind = iota
	KindKafka      Kind = iota
	KindHttp       Kind = iota
	KindGraphAPI   Kind = iota
)

var kindMap = map[Kind]string{
	KindConfig:     "Config",
	KindParse:      "Parse",
	KindFileSystem: "FileSystem",
	KindDB:         "DB",
	KindKafka:      "Kafka",
	KindHttp:       "Http",
	KindGraphAPI:   "GraphAPI",
}

func getKindTxt(kind Kind) string {
	return kindMap[kind]
}
