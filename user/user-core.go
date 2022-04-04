package user

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	. "github.com/digisan/go-generics/v2"
)

type Core struct {
	UName    string `json:"uname" validate:"required,uname"`          // unique, registered name
	Email    string `json:"email" validate:"required,email,email-db"` // unique
	Password string `json:"password" validate:"required,pwd"`         // <-- a custom validation rule, plaintext!
	key      [16]byte
}

func (c Core) String() string {
	if c.UName != "" {
		sb := strings.Builder{}
		t, v := reflect.TypeOf(c), reflect.ValueOf(c)
		for i := 0; i < t.NumField(); i++ {
			fld, val := t.Field(i), v.Field(i)
			if NotIn(fld.Name, "Password", "key") {
				sb.WriteString(fmt.Sprintf("%-12s %v\n", fld.Name+":", val.Interface()))
			}
		}
		return sb.String()
	}
	return "[Empty User]"
}

// [16]byte
func (c *Core) GenKey() [16]byte {
	if c.key == [16]byte{} {
		c.key = *(*[16]byte)([]byte(fmt.Sprintf("%d", time.Now().UnixNano())[3:19]))
	}
	return c.key
}
