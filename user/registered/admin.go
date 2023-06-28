package registered

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Admin struct {
	RegTime   time.Time `json:"regtime" validate:"regtime"`     // register time
	Active    bool      `json:"active" validate:"active"`       // true/false
	Certified bool      `json:"certified" validate:"certified"` // true/false
	Official  bool      `json:"official" validate:"official"`   // official account? true/false
	SysRole   string    `json:"role" validate:"sysRole"`        // optional
	MemLevel  uint8     `json:"level" validate:"memLevel"`      // 0-3
	MemExpire time.Time `json:"expire" validate:"memExpire"`    // optional
	Tags      string    `json:"tags" validate:"tags"`           // optional, linked by '^^'
}

func (a Admin) String() string {
	sb := strings.Builder{}
	t, v := reflect.TypeOf(a), reflect.ValueOf(a)
	for i := 0; i < t.NumField(); i++ {
		fld, val := t.Field(i), v.Field(i)
		sb.WriteString(fmt.Sprintf("%-12s %v\n", fld.Name+":", val.Interface()))
	}
	return sb.String()
}
