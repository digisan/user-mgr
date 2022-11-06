package signup

import (
	u "github.com/digisan/user-mgr/user"
)

func Store(user *u.User) error {
	user.StampRegTime()
	return u.UpdateUser(user)
}
