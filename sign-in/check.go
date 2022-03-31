package signin

import (
	"fmt"

	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
)

func CheckUserExists(login *usr.User) error {
	if udb.UserDB.UserExists(login.UName, login.Email, true) {
		return nil
	}
	if udb.UserDB.UserExists(login.UName, login.Email, false) {
		return fmt.Errorf("[%v] is dormant", login.UName)
	}
	return fmt.Errorf("[%v] is not existing", login.UName)
}

func PwdOK(login *usr.User) bool {

	mPropVal := map[string]string{
		"uname": login.UName,
		"email": login.Email,
		"phone": login.Phone,
	}

	for prop, val := range mPropVal {
		user, ok, err := udb.UserDB.LoadUserByUniProp(prop, val, true)
		if err == nil && ok && user.Password == login.Password {
			return true
		}
	}

	return false
}
