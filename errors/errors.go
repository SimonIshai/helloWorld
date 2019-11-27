package errors

import (
	errorsbasic "errors"
	"strings"
)

type errorExt struct {
	Ops  []string
	Kind Kind
	Err  error
}

func (e errorExt) Error() string {
	var ops []string
	length := len(e.Ops)
	for i := range e.Ops {
		ops = append(ops, e.Ops[length-i-1])
	}
	var msg string
	msg = getKindTxt(e.Kind)
	msg += ":"
	msg += strings.Join(e.Ops, ".")
	msg += ":"
	msg += e.Err.Error()

	return msg
}

func New(kind Kind, errMsg string) error {
	return errorExt{
		Kind: kind,
		Err:  errorsbasic.New(errMsg),
	}
}

func Wrap(err error, kind Kind, op string) error {
	if err == nil {
		return nil
	}

	switch err.(type) {

	case errorExt:
		e := err.(errorExt)
		e.Ops = append(e.Ops, op)
		return e

	default:
		return errorExt{
			Ops:  []string{op},
			Kind: kind,
			Err:  err,
		}
	}
}
