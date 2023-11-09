package resetpwd

import (
	. "github.com/digisan/user-mgr/cst"
	u "github.com/digisan/user-mgr/user"
	ur "github.com/digisan/user-mgr/user/registered"
)

// if return nil, which means user exists normally
func UserStatusIssue(login *ur.User) error {
	if u.UserExists(login.UName, login.Email, true) {
		return nil
	}
	if u.UserExists(login.UName, login.Email, false) {
		return Err(ERR_USER_DORMANT).Wrap(login.UName)
	}
	return Err(ERR_USER_NOT_EXISTS).Wrap(login.UName)
}

func EmailOK(login *ur.User) bool {
	user, ok, err := u.LoadUser(login.UName, true)
	return err == nil && ok && login.Email == user.Email
}
