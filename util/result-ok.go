package util

import "fmt"

type ResultOk struct {
	Ok  bool
	Err error
}

func NewResultOk(ok bool, failMsg string) ResultOk {
	if ok {
		return ResultOk{ok, nil}
	}
	return ResultOk{ok, fmt.Errorf("%v", failMsg)}
}
