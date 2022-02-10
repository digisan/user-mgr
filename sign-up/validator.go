package signup

import (
	"fmt"
	"unicode"

	"github.com/digisan/go-generics/str"
	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/tool"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	vf "github.com/digisan/user-mgr/user/valfield"
)

const (
	LetterLen = 6
)

var (
	condOper        = tool.CondOper
	mFieldValidator = map[string]func(string) bool{
		vf.Active:     func(v string) bool { return v == "T" || v == "F" },
		vf.UName:      func(v string) bool { return !udb.UserDB.IsExisting(v, false) },
		vf.Email:      func(v string) bool { return true },
		vf.Name:       func(v string) bool { return len(v) > 0 },
		vf.Password:   func(v string) bool { return ChkPwd(v, LetterLen) },
		vf.Avatar:     func(v string) bool { return condOper(v != "", len(v) > 10, true).(bool) },
		vf.Regtime:    func(v string) bool { return true },
		vf.Phone:      func(v string) bool { return condOper(v != "", len(v) > 6, true).(bool) },
		vf.Addr:       func(v string) bool { return condOper(v != "", len(v) > 6, true).(bool) },
		vf.SysRole:    func(v string) bool { return condOper(v != "", len(v) > 2, true).(bool) },
		vf.MemLevel:   func(v string) bool { return ChkMemLvl(v) },
		vf.MemExpire:  func(v string) bool { return condOper(v != "", len(v) > 6, true).(bool) },
		vf.NationalID: func(v string) bool { return condOper(v != "", len(v) > 6, true).(bool) },
		vf.Gender:     func(v string) bool { return condOper(v != "", v == "M" || v == "F", true).(bool) },
		vf.Position:   func(v string) bool { return condOper(v != "", len(v) > 3, true).(bool) },
		vf.Title:      func(v string) bool { return condOper(v != "", len(v) > 3, true).(bool) },
		vf.Employer:   func(v string) bool { return condOper(v != "", len(v) > 3, true).(bool) },
	}

	fEf          = fmt.Errorf
	mFieldValErr = map[string]func(string) error{
		vf.Active:     func(v string) error { return fEf("active status need 'T'/'F' for true/false") },
		vf.UName:      func(v string) error { return fEf("[%s] is already existing", v) },
		vf.Email:      func(v string) error { return fEf("invalid email format") },
		vf.Name:       func(v string) error { return fEf("invalid user real name") },
		vf.Password:   func(v string) error { return fEf("password needs minimal %d letters with UPPER,0-9,symbol", LetterLen) },
		vf.Regtime:    func(v string) error { return fEf("register time is mandatory when signing up successfully") },
		vf.Phone:      func(v string) error { return fEf("invalid telephone number") },
		vf.Addr:       func(v string) error { return fEf("invalid address") },
		vf.SysRole:    func(v string) error { return fEf("invalid system role") },
		vf.MemLevel:   func(v string) error { return fEf("invalid membership level, must between 0-9") },
		vf.MemExpire:  func(v string) error { return fEf("invalid expiry date") },
		vf.NationalID: func(v string) error { return fEf("invalid national ID") },
		vf.Gender:     func(v string) error { return fEf("gender needs 'M'/'F' for male/female") },
		vf.Position:   func(v string) error { return fEf("invalid position") },
		vf.Title:      func(v string) error { return fEf("invalid title") },
		vf.Employer:   func(v string) error { return fEf("invalid employer") },
		vf.Avatar:     func(v string) error { return fEf("invalid avatar") },
		"required":    func(v string) error { return fEf("[%s] must be provided", v) },
	}
)

func SetValidator(mExtraValidator map[string]func(string) bool) {
	// create temp mFieldValidator
	mFV := make(map[string]func(string) bool)
	for f, v := range mFieldValidator {
		mFV[f] = v
	}
	for f, v := range mExtraValidator {
		mFV[f] = v
	}
	// register
	for field, validator := range mFV {
		usr.RegisterValidator(field, validator)
	}
}

func TransInvalidErr(err error) error {
	field, tag := usr.ErrField(err)
	fn, ok := mFieldValErr[tag]
	if ok {
		return fn(tag)
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
