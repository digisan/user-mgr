package signup

import (
	"fmt"
	"time"
	"unicode"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	vf "github.com/digisan/user-mgr/user/valfield"
)

const (
	PwdLen = 6
)

var (
	fEf = fmt.Errorf

	mFieldValidator = map[string]func(any) bool{
		vf.Active:         func(v any) bool { return true },
		vf.UName:          func(v any) bool { return !udb.UserDB.UserExists(v.(string), "", false) },
		vf.Email:          func(v any) bool { return !udb.UserDB.UserExists("", v.(string), false) },
		vf.Name:           func(v any) bool { return len(v.(string)) > 0 },
		vf.Password:       func(v any) bool { return ChkPwd(v.(string), PwdLen) },
		vf.AvatarType:     func(v any) bool { return ChkAvatarType(v.(string)) },
		vf.Avatar:         func(v any) bool { return true },
		vf.Regtime:        func(v any) bool { return v != nil && v != time.Time{} },
		vf.Official:       func(v any) bool { return true },
		vf.Phone:          func(v any) bool { return v == "" || len(v.(string)) > 6 },
		vf.Country:        func(v any) bool { return v == "" || len(v.(string)) > 2 },
		vf.City:           func(v any) bool { return v == "" || len(v.(string)) > 2 },
		vf.Addr:           func(v any) bool { return v == "" || len(v.(string)) > 6 },
		vf.SysRole:        func(v any) bool { return v == "" || len(v.(string)) > 2 },
		vf.MemLevel:       func(v any) bool { return ChkMemLvl(v.(string)) },
		vf.MemExpire:      func(v any) bool { return true },
		vf.PersonalIDType: func(v any) bool { return v == "" || len(v.(string)) > 2 },
		vf.PersonalID:     func(v any) bool { return v == "" || len(v.(string)) > 6 },
		vf.Gender:         func(v any) bool { return v == "" || v == "M" || v == "F" },
		vf.DOB:            func(v any) bool { return v == "" || len(v.(string)) > 7 },
		vf.Position:       func(v any) bool { return v == "" || len(v.(string)) > 3 },
		vf.Title:          func(v any) bool { return v == "" || len(v.(string)) > 3 },
		vf.Employer:       func(v any) bool { return v == "" || len(v.(string)) > 3 },
		vf.Certified:      func(v any) bool { return true },
		vf.Bio:            func(v any) bool { return v == "" || len(v.(string)) > 3 },
		vf.Tags:           func(v any) bool { return v == "" || len(v.(string)) > 2 },
	}

	mFieldValErr = map[string]func(t, v any) error{
		vf.Active:         func(t, v any) error { return fEf("active status: true/false") },
		vf.UName:          func(t, v any) error { return fEf("[%v] is already existing", v) },
		vf.Email:          func(t, v any) error { return fEf("invalid email format OR [%v] is already registered") },
		vf.Name:           func(t, v any) error { return fEf("invalid user real name") },
		vf.Password:       func(t, v any) error { return fEf("password rule: >=%d letter with UPPER,0-9,symbol", PwdLen) },
		vf.Regtime:        func(t, v any) error { return fEf("register time is mandatory when signing up successfully") },
		vf.Official:       func(t, v any) error { return fEf("official status: true/false") },
		vf.Phone:          func(t, v any) error { return fEf("invalid telephone number") },
		vf.Country:        func(t, v any) error { return fEf("invalid country") },
		vf.City:           func(t, v any) error { return fEf("invalid city") },
		vf.Addr:           func(t, v any) error { return fEf("invalid address") },
		vf.SysRole:        func(t, v any) error { return fEf("invalid system role") },
		vf.MemLevel:       func(t, v any) error { return fEf("invalid membership level, must between 0-9") },
		vf.MemExpire:      func(t, v any) error { return fEf("invalid expiry date") },
		vf.PersonalIDType: func(t, v any) error { return fEf("invalid personal ID type") },
		vf.PersonalID:     func(t, v any) error { return fEf("invalid personal ID") },
		vf.Gender:         func(t, v any) error { return fEf("gender: 'M'/'F' for male/female") },
		vf.DOB:            func(t, v any) error { return fEf("invalid date of birth") },
		vf.Position:       func(t, v any) error { return fEf("invalid position") },
		vf.Title:          func(t, v any) error { return fEf("invalid title") },
		vf.Employer:       func(t, v any) error { return fEf("invalid employer") },
		vf.Certified:      func(t, v any) error { return fEf("certified status: true/false") },
		vf.Bio:            func(t, v any) error { return fEf("more words please") },
		vf.Tags:           func(t, v any) error { return fEf("invalid user tags") },
		vf.AvatarType:     func(t, v any) error { return fEf("invalid avatar type, must be like 'image/png'") },
		vf.Avatar:         func(t, v any) error { return fEf("invalid avatar") },
		"required":        func(t, v any) error { return fEf("[%v] is required", t) },
	}
)

func SetValidator(extraValidator map[string]func(any) bool) {
	// create temp mFieldValidator
	mFV := make(map[string]func(any) bool)
	for f, v := range mFieldValidator {
		mFV[f] = v
	}
	for f, v := range extraValidator {
		mFV[f] = v
	}
	// register
	for field, validator := range mFV {
		usr.RegisterValidator(field, validator)
	}
}

func TransInvalidErr(user *usr.User, err error) error {
	field, tag := usr.ErrField(err)
	fn, ok := mFieldValErr[tag]
	if ok {
		return fn(tag, usr.FieldValue(user, field))
	}
	lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("unknown field invalid error @ [%s]", field))
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////

func ChkPwd(s string, minLenLetter int) bool {
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
	return letters >= minLenLetter && number && upper && special
}

func ChkMemLvl(s string) bool {
	return In(s, "0", "1", "2", "3", "4", "5", "6", "7", "8", "9")
}

// <img src="data:image/png;base64,******/>
func ChkAvatarType(s string) bool {
	return s == "" || strs.HasAnyPrefix(s, "image/")
}
