package errors

type Kind int

const (
	KindDB    Kind = iota
	KindKafka Kind = iota
	KindHttp  Kind = iota
)

var kindMap = map[Kind]string{
	KindDB:    "Error.DB",
	KindKafka: "Error.Kafaka",
	KindHttp:  "Error.Http",
}

func getKindTxt(kind Kind) string {
	return kindMap[kind]
}
