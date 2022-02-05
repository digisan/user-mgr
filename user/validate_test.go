package user

import (
	"fmt"
	"testing"
)

func TestErrField(t *testing.T) {
	err := fmt.Errorf("%s", "Key: 'User.Addr' Error:Field validation for 'Addr' failed on the 'addr' tag")
	ErrField(err)
}
