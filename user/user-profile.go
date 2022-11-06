package user

import (
	"fmt"
	"reflect"
	"strings"
)

type Profile struct {
	Name           string `json:"name" validate:"name"`                     // real name
	Phone          string `json:"phone" validate:"phone,phone-db"`          //
	Country        string `json:"country" validate:"country"`               //
	City           string `json:"city" validate:"city"`                     //
	Addr           string `json:"addr" validate:"addr"`                     //
	PersonalIDType string `json:"personalidtype" validate:"personalidtype"` //
	PersonalID     string `json:"personalid" validate:"personalid"`         //
	Gender         string `json:"gender" validate:"gender"`                 //
	DOB            string `json:"dob" validate:"dob"`                       //
	Position       string `json:"position" validate:"position"`             //
	Title          string `json:"title" validate:"title"`                   //
	Employer       string `json:"employer" validate:"employer"`             //
	Bio            string `json:"bio" validate:"bio"`                       //
	AvatarType     string `json:"avatartype" validate:"avatartype"`         //
	Avatar         []byte `json:"avatar" validate:"avatar"`                 //
}

func (p Profile) String() string {
	sb := strings.Builder{}
	t, v := reflect.TypeOf(p), reflect.ValueOf(p)
	for i := 0; i < t.NumField(); i++ {
		fld, val := t.Field(i), v.Field(i)
		sb.WriteString(fmt.Sprintf("%-12s %v\n", fld.Name+":", val.Interface()))
	}
	return sb.String()
}
