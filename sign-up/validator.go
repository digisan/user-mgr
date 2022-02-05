package signup

import (
	"fmt"
	"unicode"

	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	vf "github.com/digisan/user-mgr/user/valfield"
)

func SetValidator() {
	usr.RegisterValidator(vf.Active, func(fv string) bool {
		return fv == "T" || fv == "F"
	})
	usr.RegisterValidator(vf.UName, func(fv string) bool {
		return !udb.UserDB.IsExisting(fv, false)
	})
	usr.RegisterValidator(vf.Name, func(fv string) bool {
		return len(fv) > 0
	})
	usr.RegisterValidator(vf.Password, func(fv string) bool {
		lenOK, number, upper, special := ChkPwd(fv, MinLenLetter)
		return lenOK && number && upper && special
	})
	usr.RegisterValidator(vf.Avatar, func(fv string) bool {
		return len(fv) > 0
	})
	usr.RegisterValidator(vf.Regtime, func(fv string) bool {
		return len(fv) > 0
	})
	usr.RegisterValidator(vf.Phone, func(fv string) bool {
		return len(fv) > 0
	})
	usr.RegisterValidator(vf.Addr, func(fv string) bool {
		return true
	})
	usr.RegisterValidator(vf.SysRole, func(fv string) bool {
		return true
	})
	usr.RegisterValidator(vf.SysLevel, func(fv string) bool {
		return true
	})
	usr.RegisterValidator(vf.Expire, func(fv string) bool {
		return true
	})
}

func TransInvalidErr(err error) error {
	field, tag := usr.ErrField(err)
	switch tag {
	case vf.Active:
		return fmt.Errorf("invalid active status, set 'T' for true, or 'F' for false")
	case vf.UName:
		return fmt.Errorf("invalid user name, [%s] is already existing", field)
	case vf.Email:
		return fmt.Errorf("invalid email format")
	case vf.Name:
		return fmt.Errorf("invalid user real name")
	case vf.Password:
		return fmt.Errorf("invalid password, at least %d letters, consists of UPPER-CASE, number and symbol", MinLenLetter)
	case vf.Regtime:
		return fmt.Errorf("must add register time when signing up")
	case vf.Phone:
		return fmt.Errorf("invalid telephone number")
	case vf.Addr:
		return fmt.Errorf("invalid address")
	case vf.SysRole:
		return fmt.Errorf("invalid system role")
	case vf.SysLevel:
		return fmt.Errorf("invalid system subscribe level")
	case vf.Expire:
		return fmt.Errorf("invalid expiry date")
	case vf.Avatar:
		return fmt.Errorf("invalid avatar")
	case "required":
		return fmt.Errorf("[%s] must be provided", field)
	default:
		fmt.Println(err)
		lk.FailOnErr("%v", fmt.Errorf("unknown Field Invalid Error @ [%s]", field))
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////

func ChkPwd(s string, minLenLetter int) (lenOK, number, upper, special bool) {
	letters := 0
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
	return letters >= minLenLetter, number, upper, special
}
