package udb

import (
	"fmt"
	"strings"
	"sync"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	lk "github.com/digisan/logkit"
	usr "github.com/digisan/user-mgr/user"
)

var once sync.Once

type UDB struct {
	sync.Mutex
	dbReg    *badger.DB
	dbOnline *badger.DB
}

var uDB *UDB // global, for keeping single instance

func open(dir string) *badger.DB {
	opt := badger.DefaultOptions("").WithInMemory(true)
	if dir != "" {
		opt = badger.DefaultOptions(dir)
	}
	db, err := badger.Open(opt)
	lk.FailOnErr("%v", err)
	return db
}

func getDB(dir string) *UDB {
	if uDB == nil {
		once.Do(func() {
			uDB = &UDB{
				dbReg:    open(dir),
				dbOnline: open(""),
			}
		})
	}
	return uDB
}

func (db *UDB) close() {
	db.Lock()
	defer db.Unlock()

	lk.FailOnErr("%v", db.dbOnline.Close())
	lk.FailOnErr("%v", db.dbReg.Close())
}

///////////////////////////////////////////////////////////////

func (db *UDB) LoadOnlineUser(uname string) (time.Time, error) {
	db.Lock()
	defer db.Unlock()

	tm := &time.Time{}
	err := db.dbOnline.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(uname))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return tm.UnmarshalBinary(val)
		})
	})
	return *tm, err
}

func (db *UDB) RefreshOnlineUser(uname string) error {
	db.Lock()
	defer db.Unlock()

	return db.dbOnline.Update(func(txn *badger.Txn) error {
		tmBytes, err := time.Now().UTC().MarshalBinary()
		if err != nil {
			return err
		}
		return txn.Set([]byte(uname), tmBytes)
	})
}

func (db *UDB) RemoveOnlineUser(uname string) error {
	db.Lock()
	defer db.Unlock()

	return db.dbOnline.Update(func(txn *badger.Txn) (err error) {
		return txn.Delete([]byte(uname))
	})
}

func (db *UDB) ListOnlineUsers() (unames []string, err error) {
	db.Lock()
	defer db.Unlock()

	err = db.dbOnline.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			unames = append(unames, string(it.Item().Key()))
		}
		return nil
	})
	return
}

///////////////////////////////////////////////////////////////

func (db *UDB) UpdateUser(user *usr.User) error {
	db.Lock()
	defer db.Unlock()

	// remove all existing items
	if err := db.RemoveUser(user.UName, false); err != nil {
		return err
	}
	return db.dbReg.Update(func(txn *badger.Txn) error {
		return txn.Set(user.Marshal())
	})
}

func (db *UDB) LoadUser(uname string, active bool) (*usr.User, bool, error) {
	db.Lock()
	defer db.Unlock()

	prefix := []byte("T" + usr.SEP + uname + usr.SEP)
	if !active {
		prefix = []byte("F" + usr.SEP + uname + usr.SEP)
	}

	u := &usr.User{}
	err := db.dbReg.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		if it.Seek(prefix); it.ValidForPrefix(prefix) {
			item := it.Item()
			k := item.Key()
			return item.Value(func(v []byte) error {
				// fmt.Printf("key=%s, value=%s\n", k, v)
				u.Unmarshal(k, v)
				return nil
			})
		}
		return nil
	})
	return u, u.Email != "", err
}

func (db *UDB) LoadActiveUser(uname string) (*usr.User, bool, error) {
	return db.LoadUser(uname, true)
}

func (db *UDB) LoadUserByEmail(email string, active bool) (*usr.User, bool, error) {
	db.Lock()
	defer db.Unlock()

	users, err := db.ListUsers(func(u *usr.User) bool {
		if active {
			return u.IsActive() && u.Email == email
		}
		return !u.IsActive() && u.Email == email
	})
	if len(users) > 0 {
		return users[0], true, err
	}
	return &usr.User{}, false, err
}

func (db *UDB) LoadActiveUserByEmail(email string) (*usr.User, bool, error) {
	return db.LoadUserByEmail(email, true)
}

func (db *UDB) RemoveUser(uname string, lock bool) error {
	if lock {
		db.Lock()
		defer db.Unlock()
	}

	prefixList := [][]byte{
		[]byte("T" + usr.SEP + uname + usr.SEP),
		[]byte("F" + usr.SEP + uname + usr.SEP),
	}
	return db.dbReg.Update(func(txn *badger.Txn) (err error) {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for _, prefix := range prefixList {
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				if err = txn.Delete(it.Item().KeyCopy(nil)); err != nil {
					return err
				}
			}
		}
		return err
	})
}

func (db *UDB) ListUsers(filter func(*usr.User) bool) (users []*usr.User, err error) {
	db.Lock()
	defer db.Unlock()

	err = db.dbReg.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			u := &usr.User{}
			u.Unmarshal(it.Item().Key(), nil)
			if filter(u) {
				users = append(users, u)
			}
		}
		return nil
	})
	return
}

func (db *UDB) IsExisting(uname string, onlyActive bool) bool {
	if onlyActive {
		_, okActive, err := db.LoadUser(uname, true)
		lk.WarnOnErr("%v", err)
		return okActive
	}
	_, okActive, err := db.LoadUser(uname, true)
	lk.WarnOnErr("%v", err)
	_, okDorm, err := db.LoadUser(uname, false)
	lk.WarnOnErr("%v", err)
	return okActive || okDorm
}

func (db *UDB) ActivateUser(uname string, flag bool) (bool, error) {
	u, ok, err := db.LoadUser(uname, !flag)
	if err == nil {
		if ok {
			u.Active = strings.ToUpper(fmt.Sprint(flag))[:1]
			return true, db.UpdateUser(u)
		}
		if !ok {
			return false, fmt.Errorf("no action applied for [%s]", uname)
		}
	}
	return false, err
}
