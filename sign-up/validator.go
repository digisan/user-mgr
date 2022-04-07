package signup

import (
	"fmt"
	"time"
	"unicode"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	vf "github.com/digisan/user-mgr/user/valfield"
)

var (
	fSf = fmt.Sprintf

	mFieldValidator = map[string]func(o, v any) usr.ValRst{

		vf.Active: func(o, v any) usr.ValRst {
			return usr.NewValRst(true, "")
		},

		vf.UName: func(o, v any) usr.ValRst {
			ok := !udb.UserDB.UserExists(v.(string), "", false)
			return usr.NewValRst(ok, fSf("[%v] is already existing", v))
		},

		vf.EmailDB: func(o, v any) usr.ValRst {
			ok := !udb.UserDB.UserExists("", v.(string), false)
			return usr.NewValRst(ok, fSf("[%v] is already existing", v))
		},

		vf.Name: func(o, v any) usr.ValRst {
			ok := len(v.(string)) > 0
			return usr.NewValRst(ok, "invalid user real name")
		},

		vf.Password: func(o, v any) usr.ValRst {
			return ChkPwd(v.(string))
		},

		vf.AvatarType: func(o, v any) usr.ValRst {
			return ChkAvatarType(v.(string))
		},

		vf.Avatar: func(o, v any) usr.ValRst {
			return usr.NewValRst(true, "")
		},

		vf.Regtime: func(o, v any) usr.ValRst {
			ok := v != nil && v != time.Time{}
			return usr.NewValRst(ok, "register time is mandatory when signing up successfully")
		},

		vf.Official: func(o, v any) usr.ValRst {
			return usr.NewValRst(true, "")
		},

		vf.Phone: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 6
			return usr.NewValRst(ok, "invalid telephone number")
		},

		vf.PhoneDB: func(o, v any) usr.ValRst {
			ok := v == "" || !udb.UserDB.UsedByOther(o.(*usr.User).UName, "phone", v.(string))
			return usr.NewValRst(ok, fSf("phone [%v] is already used by other user", v))
		},

		vf.Country: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return usr.NewValRst(ok, "invalid country")
		},

		vf.City: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return usr.NewValRst(ok, "invalid city")
		},

		vf.Addr: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 6
			return usr.NewValRst(ok, "invalid address")
		},

		vf.SysRole: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return usr.NewValRst(ok, "invalid system role")
		},

		vf.MemLevel: func(o, v any) usr.ValRst {
			ok := In(v.(uint8), 0, 1, 2, 3)
			return usr.NewValRst(ok, "membership level: [0-3]")
		},

		vf.MemExpire: func(o, v any) usr.ValRst {
			return usr.NewValRst(true, "")
		},

		vf.PersonalIDType: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return usr.NewValRst(ok, "invalid personal ID type")
		},

		vf.PersonalID: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 6
			return usr.NewValRst(ok, "invalid personal ID")
		},

		vf.Gender: func(o, v any) usr.ValRst {
			ok := v == "" || v == "M" || v == "F"
			return usr.NewValRst(ok, "gender: 'M'/'F' for male/female")
		},

		vf.DOB: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 7
			return usr.NewValRst(ok, "invalid date of birth")
		},

		vf.Position: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 3
			return usr.NewValRst(ok, "invalid position")
		},

		vf.Title: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 3
			return usr.NewValRst(ok, "invalid title")
		},

		vf.Employer: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return usr.NewValRst(ok, "at least 2 length for employer")
		},

		vf.Certified: func(o, v any) usr.ValRst {
			return usr.NewValRst(true, "")
		},

		vf.Bio: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 3
			return usr.NewValRst(ok, "more words please")
		},

		vf.Tags: func(o, v any) usr.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return usr.NewValRst(ok, "invalid user tags")
		},
	}
)

func SetValidator(extraValidator map[string]func(o, v any) usr.ValRst) {
	for field, validator := range MapSafeMerge(extraValidator, mFieldValidator) {
		usr.RegisterValidator(field, validator)
	}
}

////////////////////////////////////////////////////////////////////////////////////////

func ChkPwd(s string) usr.ValRst {
	pwdLen := 6
	letters, number, upper, special := 0, false, false, false
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			//return false, false, false, false
		}
	}
	ok := letters >= pwdLen && number && upper && special
	return usr.NewValRst(ok, fSf("password rule: >=%d letter with UPPER,0-9,symbol", pwdLen))
}

// <img src="data:image/png;base64,******/>
func ChkAvatarType(s string) usr.ValRst {
	ok := s == "" || strs.HasAnyPrefix(s, "image/")
	return usr.NewValRst(ok, "invalid avatar type, must be like 'image/png'")
}
