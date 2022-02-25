package signup

import (
	"fmt"
	"unicode"

	"github.com/digisan/go-generics/str"
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

	mFieldValidator = map[string]func(interface{}) bool{
		vf.Active:     func(v interface{}) bool { return v == "T" || v == "F" },
		vf.UName:      func(v interface{}) bool { return !udb.UserDB.IsExisting(v.(string), false) },
		vf.Email:      func(v interface{}) bool { return true },
		vf.Name:       func(v interface{}) bool { return len(v.(string)) > 0 },
		vf.Password:   func(v interface{}) bool { return ChkPwd(v.(string), PwdLen) },
		vf.AvatarType: func(v interface{}) bool { return ChkAvatarType(v.(string)) },
		vf.Avatar:     func(v interface{}) bool { return true },
		vf.Regtime:    func(v interface{}) bool { return true },
		vf.Phone:      func(v interface{}) bool { return v == "" || len(v.(string)) > 6 },
		vf.Addr:       func(v interface{}) bool { return v == "" || len(v.(string)) > 6 },
		vf.SysRole:    func(v interface{}) bool { return v == "" || len(v.(string)) > 2 },
		vf.MemLevel:   func(v interface{}) bool { return ChkMemLvl(v.(string)) },
		vf.MemExpire:  func(v interface{}) bool { return v == "" || len(v.(string)) > 6 },
		vf.NationalID: func(v interface{}) bool { return v == "" || len(v.(string)) > 6 },
		vf.Gender:     func(v interface{}) bool { return v == "" || v == "M" || v == "F" },
		vf.Position:   func(v interface{}) bool { return v == "" || len(v.(string)) > 3 },
		vf.Title:      func(v interface{}) bool { return v == "" || len(v.(string)) > 3 },
		vf.Employer:   func(v interface{}) bool { return v == "" || len(v.(string)) > 3 },
		vf.Tags:       func(v interface{}) bool { return v == "" || len(v.(string)) > 2 },
	}

	mFieldValErr = map[string]func(t, v interface{}) error{
		vf.Active:     func(t, v interface{}) error { return fEf("active status: 'T'/'F' for true/false") },
		vf.UName:      func(t, v interface{}) error { return fEf("[%v] is already existing", v) },
		vf.Email:      func(t, v interface{}) error { return fEf("invalid email format") },
		vf.Name:       func(t, v interface{}) error { return fEf("invalid user real name") },
		vf.Password:   func(t, v interface{}) error { return fEf("password rule: >=%d letter with UPPER,0-9,symbol", PwdLen) },
		vf.Regtime:    func(t, v interface{}) error { return fEf("register time is mandatory when signing up successfully") },
		vf.Phone:      func(t, v interface{}) error { return fEf("invalid telephone number") },
		vf.Addr:       func(t, v interface{}) error { return fEf("invalid address") },
		vf.SysRole:    func(t, v interface{}) error { return fEf("invalid system role") },
		vf.MemLevel:   func(t, v interface{}) error { return fEf("invalid membership level, must between 0-9") },
		vf.MemExpire:  func(t, v interface{}) error { return fEf("invalid expiry date") },
		vf.NationalID: func(t, v interface{}) error { return fEf("invalid national ID") },
		vf.Gender:     func(t, v interface{}) error { return fEf("gender: 'M'/'F' for male/female") },
		vf.Position:   func(t, v interface{}) error { return fEf("invalid position") },
		vf.Title:      func(t, v interface{}) error { return fEf("invalid title") },
		vf.Employer:   func(t, v interface{}) error { return fEf("invalid employer") },
		vf.Tags:       func(t, v interface{}) error { return fEf("invalid user tags") },
		vf.AvatarType: func(t, v interface{}) error { return fEf("invalid avatar type, must be like image/png") },
		vf.Avatar:     func(t, v interface{}) error { return fEf("invalid avatar") },
		"required":    func(t, v interface{}) error { return fEf("[%v] is required", t) },
	}
)

func SetValidator(extraValidator map[string]func(interface{}) bool) {
	// create temp mFieldValidator
	mFV := make(map[string]func(interface{}) bool)
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
		return fn(tag, user.FieldValue(field))
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
	return str.In(s, "0", "1", "2", "3", "4", "5", "6", "7", "8", "9")
}

// <img src="data:image/png;base64,******/>
func ChkAvatarType(s string) bool {
	return s == "" || strs.HasAnyPrefix(s, "image/")
}
