package signup

import (
	"fmt"
	"time"
	"unicode"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	u "github.com/digisan/user-mgr/user"
	vf "github.com/digisan/user-mgr/user/valfield"
)

const (
	minPwdLen = 8
)

var (
	fSf = fmt.Sprintf

	mFieldValidator = map[string]func(o, v any) u.ValRst{

		vf.Active: func(o, v any) u.ValRst {
			return u.NewValRst(true, "")
		},

		vf.UName: func(o, v any) u.ValRst {
			return ChkUName(v.(string))
		},

		vf.EmailDB: func(o, v any) u.ValRst {
			ok := !u.UserExists("", v.(string), false)
			return u.NewValRst(ok, fSf("[%v] is already existing", v))
		},

		vf.Name: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValRst(ok, "invalid user real name")
		},

		vf.Password: func(o, v any) u.ValRst {
			return ChkPwd(v.(string))
		},

		vf.AvatarType: func(o, v any) u.ValRst {
			return ChkAvatarType(v.(string))
		},

		vf.Avatar: func(o, v any) u.ValRst {
			return u.NewValRst(true, "")
		},

		vf.RegTime: func(o, v any) u.ValRst {
			ok := v != nil && v != time.Time{}
			return u.NewValRst(ok, "register time is mandatory when signing up successfully")
		},

		vf.Official: func(o, v any) u.ValRst {
			return u.NewValRst(true, "")
		},

		vf.Phone: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 6
			return u.NewValRst(ok, "invalid telephone number")
		},

		vf.PhoneDB: func(o, v any) u.ValRst {
			ok := v == "" || !u.UsedByOther(o.(*u.User).UName, "phone", v.(string))
			return u.NewValRst(ok, fSf("phone [%v] is already used by other user", v))
		},

		vf.Country: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValRst(ok, "invalid country")
		},

		vf.City: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValRst(ok, "invalid city")
		},

		vf.Addr: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 6
			return u.NewValRst(ok, "invalid address")
		},

		vf.SysRole: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValRst(ok, "invalid system role")
		},

		vf.MemLevel: func(o, v any) u.ValRst {
			ok := In(v.(uint8), 0, 1, 2, 3)
			return u.NewValRst(ok, "membership level: [0-3]")
		},

		vf.MemExpire: func(o, v any) u.ValRst {
			return u.NewValRst(true, "")
		},

		vf.PersonalIDType: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValRst(ok, "invalid personal ID type")
		},

		vf.PersonalID: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 6
			return u.NewValRst(ok, "invalid personal ID")
		},

		vf.Gender: func(o, v any) u.ValRst {
			ok := v == "" || v == "M" || v == "F"
			return u.NewValRst(ok, "gender: 'M'/'F' for male/female")
		},

		vf.DOB: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 7
			return u.NewValRst(ok, "invalid date of birth")
		},

		vf.Position: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 3
			return u.NewValRst(ok, "invalid position")
		},

		vf.Title: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 3
			return u.NewValRst(ok, "invalid title")
		},

		vf.Employer: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValRst(ok, "at least 2 length for employer")
		},

		vf.Certified: func(o, v any) u.ValRst {
			return u.NewValRst(true, "")
		},

		vf.Bio: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 3
			return u.NewValRst(ok, "more words please")
		},

		vf.Tags: func(o, v any) u.ValRst {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValRst(ok, "invalid user tags")
		},
	}
)

func SetValidator(extraValidator map[string]func(o, v any) u.ValRst) {
	for field, validator := range MapSafeMerge(extraValidator, mFieldValidator) {
		u.RegisterValidator(field, validator)
	}
}

////////////////////////////////////////////////////////////////////////////////////////

func ChkUName(s string) u.ValRst {
	for _, c := range s {
		if unicode.IsPunct(c) || unicode.IsSymbol(c) || unicode.IsSpace(c) {
			if NotIn(c, '.', '-', '_') { // only allow '.' '-' '_' in user name
				return u.NewValRst(false, fSf("user name: [%v] has invalid character [%v]", s, string(c)))
			}
		}
	}
	ok := !u.UserExists(s, "", false)
	return u.NewValRst(ok, fSf("[%v] is already existing", s))
}

func ChkPwd(s string) u.ValRst {
	number, lower, upper, special := false, false, false, false
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case c == '	' || c == '\t':
			return u.NewValRst(false, "[space] and [table] are not allowed in password")

		// case unicode.IsLetter(c) || c == ' ':

		default:
			//return false, false, false, false
		}
	}
	ok := len(s) >= minPwdLen && (number && lower && upper && special)
	return u.NewValRst(ok, PwdRule())
}

func PwdRule() string {
	return fSf("No less than %d Characters including at least one UPPER, Number and Symbol", minPwdLen)
}

// <img src="data:image/png;base64,******/>
func ChkAvatarType(s string) u.ValRst {
	ok := s == "" || strs.HasAnyPrefix(s, "image/")
	return u.NewValRst(ok, "invalid avatar type, must be like 'image/png'")
}
