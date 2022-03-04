package resetpwd

import (
	"fmt"

	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
)

func UserExists(login *usr.User) error {
	if udb.UserDB.IsExisting(login.UName, login.Email, true) {
		return nil
	}
	if udb.UserDB.IsExisting(login.UName, login.Email, false) {
		return fmt.Errorf("[%v] is dormant", login.UName)
	}
	return fmt.Errorf("[%v] is not existing", login.UName)
}

func EmailOK(login *usr.User) bool {
	user, ok, err := udb.UserDB.LoadUser(login.UName, true)
	return err == nil && ok && login.Email == user.Email
}
