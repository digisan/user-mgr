package signin

import (
	"fmt"

	u "github.com/digisan/user-mgr/user"
)

func CheckUserExists(login *u.User) error {
	if u.UserExists(login.UName, login.Email, true) {
		return nil
	}
	if u.UserExists(login.UName, login.Email, false) {
		return fmt.Errorf("[%v] is dormant", login.UName)
	}
	return fmt.Errorf("[%v] is not existing", login.UName)
}

// if successful, then update login user
func PwdOK(login *u.User) bool {

	mPropVal := map[string]string{
		"uname": login.UName,
		"email": login.Email,
		"phone": login.Phone,
	}

	for prop, val := range mPropVal {
		if len(val) == 0 {
			continue
		}
		user, ok, err := u.LoadUserByUniProp(prop, val, true)
		if err == nil && ok && user.Password == login.Password {
			*login = *user
			return true
		}
	}

	return false
}
