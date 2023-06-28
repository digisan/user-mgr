package user

import (
	bh "github.com/digisan/db-helper/badger"
	"github.com/digisan/user-mgr/user/db"
	. "github.com/digisan/user-mgr/user/online"
)

///////////////////////////////////////////////////

func GetOnline(uname string) (*User, error) {
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	return bh.GetOneObject[User]([]byte(uname))
}

func RefreshOnline(uname string) (*User, error) {
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	u := NewUser(uname)
	return u, bh.UpsertOneObject(u)
}

func RmOnline(uname string) (int, error) {
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	return bh.DeleteOneObject[User]([]byte(uname))
}

func OnlineUsers() ([]*User, error) {
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	return bh.GetObjects[User]([]byte(""), nil)
}
