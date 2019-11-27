package errors

type Kind int

const (
	KindParse      Kind = iota
	KindFileSystem Kind = iota
	KindDB         Kind = iota
	KindKafka      Kind = iota
	KindHttp       Kind = iota
)

var kindMap = map[Kind]string{
	KindParse:      "Error.Parse",
	KindFileSystem: "Error.FileSystem",
	KindDB:         "Error.DB",
	KindKafka:      "Error.Kafka",
	KindHttp:       "Error.Http",
}

func getKindTxt(kind Kind) string {
	return kindMap[kind]
}
