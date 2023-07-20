package signup

import (
	"fmt"
	"time"
	"unicode"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	u "github.com/digisan/user-mgr/user"
	ur "github.com/digisan/user-mgr/user/registered"
)

const (
	minPwdLen = 8
)

var (
	fSf = fmt.Sprintf

	mFieldValidator = map[string]func(o, v any) u.ValidateResult{

		u.Active: func(o, v any) u.ValidateResult {
			return u.NewValidateResult(true, "")
		},

		u.UName: func(o, v any) u.ValidateResult {
			return CheckUName(v.(string))
		},

		u.EmailDB: func(o, v any) u.ValidateResult {
			ok := !u.UserExists("", v.(string), false)
			return u.NewValidateResult(ok, fSf("[%v] is already existing", v))
		},

		u.Name: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValidateResult(ok, "invalid user real name")
		},

		u.Password: func(o, v any) u.ValidateResult {
			return CheckPwd(v.(string))
		},

		u.AvatarType: func(o, v any) u.ValidateResult {
			return CheckAvatarType(v.(string))
		},

		u.Avatar: func(o, v any) u.ValidateResult {
			return u.NewValidateResult(true, "")
		},

		u.RegTime: func(o, v any) u.ValidateResult {
			ok := v != nil && v != time.Time{}
			return u.NewValidateResult(ok, "register time is mandatory when signing up successfully")
		},

		u.Official: func(o, v any) u.ValidateResult {
			return u.NewValidateResult(true, "")
		},

		u.Phone: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 6
			return u.NewValidateResult(ok, "invalid telephone number")
		},

		u.PhoneDB: func(o, v any) u.ValidateResult {
			ok := v == "" || !u.UsedByOther(o.(*ur.User).UName, "phone", v.(string))
			return u.NewValidateResult(ok, fSf("phone [%v] is already used by other user", v))
		},

		u.Country: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValidateResult(ok, "invalid country")
		},

		u.City: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValidateResult(ok, "invalid city")
		},

		u.Addr: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 6
			return u.NewValidateResult(ok, "invalid address")
		},

		u.SysRole: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValidateResult(ok, "invalid system role")
		},

		u.MemLevel: func(o, v any) u.ValidateResult {
			ok := In(v.(uint8), 0, 1, 2, 3)
			return u.NewValidateResult(ok, "membership level: [0-3]")
		},

		u.MemExpire: func(o, v any) u.ValidateResult {
			return u.NewValidateResult(true, "")
		},

		u.PersonalIDType: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValidateResult(ok, "invalid personal ID type")
		},

		u.PersonalID: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 6
			return u.NewValidateResult(ok, "invalid personal ID")
		},

		u.Gender: func(o, v any) u.ValidateResult {
			ok := v == "" || v == "M" || v == "F"
			return u.NewValidateResult(ok, "gender: 'M'/'F' for male/female")
		},

		u.DOB: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 7
			return u.NewValidateResult(ok, "invalid date of birth")
		},

		u.Position: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 3
			return u.NewValidateResult(ok, "invalid position")
		},

		u.Title: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 3
			return u.NewValidateResult(ok, "invalid title")
		},

		u.Employer: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValidateResult(ok, "at least 2 length for employer")
		},

		u.Certified: func(o, v any) u.ValidateResult {
			return u.NewValidateResult(true, "")
		},

		u.Bio: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 3
			return u.NewValidateResult(ok, "more words please")
		},

		u.Notes: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValidateResult(ok, "invalid user notes")
		},

		u.Status: func(o, v any) u.ValidateResult {
			ok := v == "" || len(v.(string)) > 2
			return u.NewValidateResult(ok, "invalid user status")
		},
	}
)

func SetValidator(extraValidator map[string]func(o, v any) u.ValidateResult) {
	for field, validator := range MapSafeMerge(extraValidator, mFieldValidator) {
		u.RegisterValidator(field, validator)
	}
}

////////////////////////////////////////////////////////////////////////////////////////

func CheckUName(s string) u.ValidateResult {
	for _, c := range s {
		if unicode.IsPunct(c) || unicode.IsSymbol(c) || unicode.IsSpace(c) {
			if NotIn(c, '.', '-', '_') { // only allow '.' '-' '_' in user name
				return u.NewValidateResult(false, fSf("user name: [%v] has invalid character [%v]", s, string(c)))
			}
		}
	}
	ok := !u.UserExists(s, "", false)
	return u.NewValidateResult(ok, fSf("[%v] is already existing", s))
}

func CheckPwd(s string) u.ValidateResult {
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
			return u.NewValidateResult(false, "[space] and [table] are not allowed in password")

		// case unicode.IsLetter(c) || c == ' ':

		default:
			//return false, false, false, false
		}
	}
	ok := len(s) >= minPwdLen && (number && lower && upper && special)
	return u.NewValidateResult(ok, "Password Rule: "+PwdRule())
}

func PwdRule() string {
	return fSf("No less than %d Characters including at least one UPPER, Number and Symbol", minPwdLen)
}

// <img src="data:image/png;base64,******/>
func CheckAvatarType(s string) u.ValidateResult {
	ok := s == "" || strs.HasAnyPrefix(s, "image/")
	return u.NewValidateResult(ok, "invalid avatar type, must be like 'image/png'")
}
