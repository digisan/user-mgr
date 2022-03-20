package user

import (
	"fmt"
	"reflect"
	"strings"
)

type Profile struct {
	Name       string `json:"name" validate:"required,name"`    // real name
	Phone      string `json:"phone" validate:"phone"`           // optional
	Country    string `json:"country" validate:"country"`       // optional
	City       string `json:"city" validate:"city"`             // optional
	Addr       string `json:"addr" validate:"addr"`             // optional
	NationalID string `json:"nationalid" validate:"nationalid"` // optional
	Gender     string `json:"gender" validate:"gender"`         // optional
	DOB        string `json:"dob" validate:"dob"`               // optional
	Position   string `json:"position" validate:"position"`     // optional
	Title      string `json:"title" validate:"title"`           // optional
	Employer   string `json:"employer" validate:"employer"`     // optional
	Bio        string `json:"bio" validate:"bio"`               // optional
	AvatarType string `json:"avatartype" validate:"avatartype"` // optional
	Avatar     []byte `json:"avatar" validate:"avatar"`         // optional
}

func (p Profile) String() string {
	sb := strings.Builder{}
	t, v := reflect.TypeOf(p), reflect.ValueOf(p)
	for i := 0; i < t.NumField(); i++ {
		fld, val := t.Field(i), v.Field(i)
		sb.WriteString(fmt.Sprintf("%-12s %v\n", fld.Name+":", val.String()))
	}
	return sb.String()
}
