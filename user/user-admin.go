package user

import (
	"fmt"
	"reflect"
	"strings"
)

type Admin struct {
	Regtime   string `json:"regtime" validate:"regtime"`   // register time
	Active    string `json:"active" validate:"active"`     // "T" "F"
	SysRole   string `json:"role" validate:"sysRole"`      // optional
	MemLevel  string `json:"level" validate:"memLevel"`    // optional
	MemExpire string `json:"expire" validate:"memExpire"`  // optional
	Official  string `json:"official" validate:"official"` // official account? "T" "F"
	Tags      string `json:"tags" validate:"tags"`         // optional // linked by '^'
}

func (a Admin) String() string {
	sb := strings.Builder{}
	t, v := reflect.TypeOf(a), reflect.ValueOf(a)
	for i := 0; i < t.NumField(); i++ {
		fld, val := t.Field(i), v.Field(i)
		sb.WriteString(fmt.Sprintf("%-12s %v\n", fld.Name+":", val.String()))
	}
	return sb.String()
}
