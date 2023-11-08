package signin

import (
	"fmt"

	u "github.com/digisan/user-mgr/user"
	ur "github.com/digisan/user-mgr/user/registered"
)

// if return nil, which means user exists normally
func UserStatusIssue(login *ur.User) error {
	uname, email := login.UName, login.Email
	if u.UserExists(uname, email, true) {
		return nil
	}
	if u.UserExists(uname, email, false) {
		return fmt.Errorf("[%v] is dormant, cannot login", uname)
	}
	return fmt.Errorf("[%v] is not existing", uname)
}

// if successful, then update login user
func PwdOK(login *ur.User) bool {
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
