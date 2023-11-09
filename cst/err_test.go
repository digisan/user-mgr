package cst

import (
	"errors"
	"fmt"
	"testing"
)

func TestErr(t *testing.T) {
	{
		err := NewErr(ERR_USER_INV_FIELD)
		fmt.Println(err, err.Code())
		fmt.Println(GetErrCode(err), ":", GetErrMsg(GetErrCode(err)))
	}
	{
		err := NewErr(ERR_TIMEOUT).Wrap(123)
		fmt.Println(err)
		fmt.Println(GetErrCode(err), ":", GetErrMsg(GetErrCode(err)))
	}
	{
		err := errors.New("sub error")
		err = NewErr(ERR_USER_DORMANT).Wrap(err)
		fmt.Println(err)
		fmt.Println(GetErrCode(err), ":", GetErrMsg(GetErrCode(err)))
	}
}
