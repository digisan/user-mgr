package signout

import "github.com/digisan/user-mgr/udb"

func Logout(uname string) error {
	return udb.UserDB.RemoveOnlineUser(uname)
}
