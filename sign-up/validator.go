package signup

import (
	"fmt"
	"time"
	"unicode"

	. "github.com/digisan/go-generics"
	"github.com/digisan/gotk/strs"
	. "github.com/digisan/user-mgr/cst"
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
			return NewResultOk(true, nil)
		},

		ur.UName: func(o, v any) ResultOk {
			return CheckUName(v.(string))
		},

		ur.EmailDB: func(o, v any) ResultOk {
			ok := !UserExists("", v.(string), false)
			return NewResultOk(ok, Err(ERR_USER_ALREADY_REG).Wrap(v))
		},

		ur.Name: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("real name"))
		},

		ur.Password: func(o, v any) ResultOk {
			return CheckPwd(v.(string))
		},

		ur.AvatarType: func(o, v any) ResultOk {
			return CheckAvatarType(v.(string))
		},

		ur.Avatar: func(o, v any) ResultOk {
			return NewResultOk(true, nil)
		},

		ur.RegTime: func(o, v any) ResultOk {
			ok := v != nil && v != time.Time{}
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("register timestamp is mandatory if sign up successfully"))
		},

		ur.Official: func(o, v any) ResultOk {
			return NewResultOk(true, nil)
		},

		ur.Phone: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 6
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("phone number"))
		},

		ur.PhoneDB: func(o, v any) ResultOk {
			ok := v == "" || !UsedByOther(o.(*ur.User).UName, "phone", v.(string))
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap(fmt.Sprintf("phone [%v] is occupied", v)))
		},

		ur.Country: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("country"))
		},

		ur.City: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("city"))
		},

		ur.Addr: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 6
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("address"))
		},

		ur.SysRole: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("system role"))
		},

		ur.MemLevel: func(o, v any) ResultOk {
			ok := In(v.(uint8), 0, 1, 2, 3)
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("membership level: [0-3]"))
		},

		ur.MemExpire: func(o, v any) ResultOk {
			return NewResultOk(true, nil)
		},

		ur.PersonalIDType: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("personal ID type"))
		},

		ur.PersonalID: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 6
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("personal ID"))
		},

		ur.Gender: func(o, v any) ResultOk {
			ok := v == "" || v == "M" || v == "F"
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("gender: 'M'/'F'"))
		},

		ur.DOB: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 7
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("date of birth"))
		},

		ur.Position: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 3
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("position"))
		},

		ur.Title: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 3
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("title"))
		},

		ur.Employer: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("employer"))
		},

		ur.Certified: func(o, v any) ResultOk {
			return NewResultOk(true, nil)
		},

		ur.Bio: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("bio"))
		},

		ur.Notes: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("notes"))
		},

		ur.Status: func(o, v any) ResultOk {
			ok := v == "" || len(v.(string)) > 2
			return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("status"))
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
				return NewResultOk(false, Err(ERR_USER_INV_FIELD).Wrap("user name (only allow '.' '-' '_')"))
			}
		}
	}
	ok := !UserExists(s, "", false)
	return NewResultOk(ok, Err(ERR_USER_ALREADY_REG).Wrap(s))
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
			return NewResultOk(false, Err(ERR_USER_INV_FIELD).Wrap("password rule: [blank] is not allowed"))

		// case unicode.IsLetter(c) || c == ' ':

		default:
			//return false, false, false, false
		}
	}
	ok := len(s) >= minPwdLen && (number && lower && upper && special)
	return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("password rule"))
}

func PwdRule() string {
	return fSf("No less than %d Characters including at least one UPPER, Number and Symbol", minPwdLen)
}

// <img src="data:image/png;base64,******/>
func CheckAvatarType(s string) ResultOk {
	ok := s == "" || strs.HasAnyPrefix(s, "image/")
	return NewResultOk(ok, Err(ERR_USER_INV_FIELD).Wrap("avatar type, must be like 'image/png'"))
}
