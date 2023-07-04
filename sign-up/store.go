package signup

import (
	u "github.com/digisan/user-mgr/user"
	ur "github.com/digisan/user-mgr/user/registered"
)

func Store(user *ur.User) error {
	user.StampRegTime()
	return u.UpdateUser(user)
}
