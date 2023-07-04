package resetpwd

import (
	"fmt"

	u "github.com/digisan/user-mgr/user"
	ur "github.com/digisan/user-mgr/user/registered"
)

// if return nil, which means user exists normally
func UserStatusIssue(login *ur.User) error {
	if u.UserExists(login.UName, login.Email, true) {
		return nil
	}
	if u.UserExists(login.UName, login.Email, false) {
		return fmt.Errorf("[%v] is dormant", login.UName)
	}
	return fmt.Errorf("[%v] is not existing", login.UName)
}

func EmailOK(login *ur.User) bool {
	user, ok, err := u.LoadUser(login.UName, true)
	return err == nil && ok && login.Email == user.Email
}
