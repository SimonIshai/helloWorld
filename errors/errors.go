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
	msg += `{"error":`
	msg += "{\"kind\":\"" + getKindTxt(e.Kind) + "\""
	msg += `, "trace":`
	msg += `["` + strings.Join(ops, `","`) + "\"]"
	msg += `, "err_msg":"` + e.Err.Error() + "\"}"
	msg += `}`
	return msg
}

func New(kind Kind, errMsg string) error {
	return errorExt{
		Kind: kind,
		Err:  errorsbasic.New(errMsg),
	}
}

func Wrap(err error, op string) error {
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
			Ops: []string{op},
			Err: err,
		}
	}
}

func WrapWithKind(err error, kind Kind, op string) error {
	if err == nil {
		return nil
	}

	switch err.(type) {

	case errorExt:
		e := err.(errorExt)
		e.Ops = append(e.Ops, op)
		e.Kind = kind
		return e

	default:
		return errorExt{
			Ops:  []string{op},
			Kind: kind,
			Err:  err,
		}
	}
}

// Wrap accepts a set of arguments.
// Example: Wrap(err, KindHttp, "api call")
func Wrap2(args ...interface{}) error { //err error, kind Kind, op string) error {
	var err error
	var kind Kind
	var op string

	for i := range args {
		switch args[i].(type) {
		case error:
			err = args[i].(error)
			break
		case Kind:
			kind = args[i].(Kind)
			break
		case string:
			op = args[i].(string)
		}
	}
	return WrapWithKind(err, kind, op)
}
