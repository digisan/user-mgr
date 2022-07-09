package signout

import u "github.com/digisan/user-mgr/user"

func Logout(uname string) error {
	_, err := u.RmOnline(uname)
	return err
}
