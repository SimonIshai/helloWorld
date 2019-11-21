package errors

type ErrorExt struct {
	op  string
	msg string
}

func (e *ErrorExt) Error() string {
	return e.msg
}
