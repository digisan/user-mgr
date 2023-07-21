package signup

import (
	"fmt"
	"time"
	"unicode"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	. "github.com/digisan/user-mgr/user"
	ur "github.com/digisan/user-mgr/user/registered"
)

const (
	minPwdLen = 8
)

var (
	fSf = fmt.Sprintf

	mFieldValidator = map[string]func(o, v any) ResultOk{

		ur.Active: func(o, v any) ResultOk {
			return NewResultOk(true, "")
		},

		ur.UName: func(o, v any) ResultOk {
			return CheckUName(v.(string))
		},

		ur.EmailDB: func(o, v any) ResultOk {
			ok := !UserExists("", v.(string), false)
			return NewResultOk(ok, fSf("[%v] is already existing", v))
		},

		ur.Name: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, "invalid user real name")
		},

		ur.Password: func(o, v any) ResultOk {
			return CheckPwd(v.(string))
		},

		ur.AvatarType: func(o, v any) ResultOk {
			return CheckAvatarType(v.(string))
		},

		ur.Avatar: func(o, v any) ResultOk {
			return NewResultOk(true, "")
		},

		ur.RegTime: func(o, v any) ResultOk {
			ok := v != nil && v != time.Time{}
			return NewResultOk(ok, "register time is mandatory when signing up successfully")
		},

		ur.Official: func(o, v any) ResultOk {
			return NewResultOk(true, "")
		},

		ur.Phone: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 6
			return NewResultOk(ok, "invalid telephone number")
		},

		ur.PhoneDB: func(o, v any) ResultOk {
			ok := v == "" || !UsedByOther(o.(*ur.User).UName, "phone", v.(string))
			return NewResultOk(ok, fSf("phone [%v] is used by others", v))
		},

		ur.Country: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, "invalid country")
		},

		ur.City: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, "invalid city")
		},

		ur.Addr: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 6
			return NewResultOk(ok, "invalid address")
		},

		ur.SysRole: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, "invalid system role")
		},

		ur.MemLevel: func(o, v any) ResultOk {
			ok := In(v.(uint8), 0, 1, 2, 3)
			return NewResultOk(ok, "membership level: [0-3]")
		},

		ur.MemExpire: func(o, v any) ResultOk {
			return NewResultOk(true, "")
		},

		ur.PersonalIDType: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, "invalid personal ID type")
		},

		ur.PersonalID: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 6
			return NewResultOk(ok, "invalid personal ID")
		},

		ur.Gender: func(o, v any) ResultOk {
			ok := v == "" || v == "M" || v == "F"
			return NewResultOk(ok, "gender: 'M'/'F' for male/female")
		},

		ur.DOB: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 7
			return NewResultOk(ok, "invalid date of birth")
		},

		ur.Position: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 3
			return NewResultOk(ok, "invalid position")
		},

		ur.Title: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 3
			return NewResultOk(ok, "invalid title")
		},

		ur.Employer: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, "at least 2 length for employer")
		},

		ur.Certified: func(o, v any) ResultOk {
			return NewResultOk(true, "")
		},

		ur.Bio: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 3
			return NewResultOk(ok, "more words please")
		},

		ur.Notes: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, "invalid user notes")
		},

		ur.Status: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, "invalid user status")
		},
	}
)

func SetValidator(extraValidator map[string]func(o, v any) ResultOk) {
	for field, validator := range MapSafeMerge(extraValidator, mFieldValidator) {
		RegisterValidator(field, validator)
	}
}

////////////////////////////////////////////////////////////////////////////////////////

func CheckUName(s string) ResultOk {
	for _, c := range s {
		if unicode.IsPunct(c) || unicode.IsSymbol(c) || unicode.IsSpace(c) {
			if NotIn(c, '.', '-', '_') { // only allow '.' '-' '_' in user name
				return NewResultOk(false, fSf("user name: [%v] has invalid character [%v]", s, string(c)))
			}
		}
	}
	ok := !UserExists(s, "", false)
	return NewResultOk(ok, fSf("[%v] is already existing", s))
}

func CheckPwd(s string) ResultOk {
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
			return NewResultOk(false, "[space] and [table] are not allowed in password")

		// case unicode.IsLetter(c) || c == ' ':

		default:
			//return false, false, false, false
		}
	}
	ok := len(s) >= minPwdLen && (number && lower && upper && special)
	return NewResultOk(ok, "Password Rule: "+PwdRule())
}

func PwdRule() string {
	return fSf("No less than %d Characters including at least one UPPER, Number and Symbol", minPwdLen)
}

// <img src="data:image/png;base64,******/>
func CheckAvatarType(s string) ResultOk {
	ok := s == "" || strs.HasAnyPrefix(s, "image/")
	return NewResultOk(ok, "invalid avatar type, must be like 'image/png'")
}
