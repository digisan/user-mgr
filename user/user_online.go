package user

import (
	"fmt"

	bh "github.com/digisan/db-helper/badger"
	"github.com/digisan/user-mgr/db"
	. "github.com/digisan/user-mgr/user/online"
)

func GetOnline(uname string) (*User, error) {
	if !db.IsInit() {
		return nil, fmt.Errorf("db is not initialized")
	}
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	return bh.GetOneObject[User]([]byte(uname))
}

func RefreshOnline(uname string) (*User, error) {
	if !db.IsInit() {
		return nil, fmt.Errorf("db is not initialized")
	}
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	u := NewUser(uname)
	return u, bh.UpsertOneObject(u)
}

func RmOnline(uname string) (int, error) {
	if !db.IsInit() {
		return -1, fmt.Errorf("db is not initialized")
	}
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	return bh.DeleteOneObject[User]([]byte(uname))
}

func OnlineUsers() ([]*User, error) {
	if !db.IsInit() {
		return nil, fmt.Errorf("db is not initialized")
	}
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	return bh.GetObjects[User]([]byte(""), nil)
}
