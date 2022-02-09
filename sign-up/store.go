package signup

import (
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
)

func Store(user usr.User) error {
	user.StampRegTime()
	return udb.UserDB.UpdateUser(user)
}
