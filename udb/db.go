package udb

import (
	"sync"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	lk "github.com/digisan/logkit"
	usr "github.com/digisan/user-mgr/user"
)

var once sync.Once

type UDB struct {
	sync.Mutex
	db       *badger.DB
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

func GetDB(dir string) *UDB {
	if uDB == nil {
		once.Do(func() {
			uDB = &UDB{
				db:       open(dir),
				dbOnline: open(""),
			}
		})
	}
	return uDB
}

func (udb *UDB) Close() {
	udb.Lock()
	defer udb.Unlock()

	lk.FailOnErr("%v", udb.dbOnline.Close())
	lk.FailOnErr("%v", udb.db.Close())
}

///////////////////////////////////////////////////////////////

func (udb *UDB) LoadOnlineUser(uname string) (time.Time, error) {
	udb.Lock()
	defer udb.Unlock()

	tm := &time.Time{}
	err := udb.dbOnline.View(func(txn *badger.Txn) error {
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

func (udb *UDB) UpdateOnlineUser(uname string) error {
	udb.Lock()
	defer udb.Unlock()

	return udb.dbOnline.Update(func(txn *badger.Txn) error {
		tmBytes, err := time.Now().UTC().MarshalBinary()
		if err != nil {
			return err
		}
		return txn.Set([]byte(uname), tmBytes)
	})
}

func (udb *UDB) RemoveOnlineUser(uname string) error {
	udb.Lock()
	defer udb.Unlock()

	return udb.dbOnline.Update(func(txn *badger.Txn) (err error) {
		return txn.Delete([]byte(uname))
	})
}

func (udb *UDB) ListOnlineUsers() (unames []string, err error) {
	udb.Lock()
	defer udb.Unlock()

	err = udb.dbOnline.View(func(txn *badger.Txn) error {
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

func (udb *UDB) UpdateUser(user usr.User) error {
	// remove all existing items
	udb.RemoveUser(user.UName)

	udb.Lock()
	defer udb.Unlock()
	return udb.db.Update(func(txn *badger.Txn) error {
		return txn.Set(user.Marshal())
	})
}

func (udb *UDB) LoadUser(uname string, active bool) (usr.User, bool, error) {
	udb.Lock()
	defer udb.Unlock()

	u := &usr.User{}
	err := udb.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("T||" + uname + "||")
		if !active {
			prefix = []byte("F||" + uname + "||")
		}
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
	return *u, u.Email != "", err
}

func (udb *UDB) LoadActiveUser(uname string) (usr.User, bool, error) {
	return udb.LoadUser(uname, true)
}

func (udb *UDB) LoadUserByEmail(email string, active bool) (usr.User, bool, error) {
	udb.Lock()
	defer udb.Unlock()

	users, err := udb.ListUsers(func(u *usr.User) bool {
		if active {
			return u.IsActive() && u.Email == email
		}
		return !u.IsActive() && u.Email == email
	})
	if len(users) > 0 {
		return users[0], true, err
	}
	return usr.User{}, false, err
}

func (udb *UDB) LoadActiveUserByEmail(email string) (usr.User, bool, error) {
	return udb.LoadUserByEmail(email, true)
}

func (udb *UDB) RemoveUser(uname string) error {
	udb.Lock()
	defer udb.Unlock()

	return udb.db.Update(func(txn *badger.Txn) (err error) {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefixList := [][]byte{
			[]byte("T||" + uname + "||"),
			[]byte("F||" + uname + "||"),
		}
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

func (udb *UDB) ListUsers(filter func(*usr.User) bool) (users []usr.User, err error) {
	udb.Lock()
	defer udb.Unlock()

	err = udb.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			u := &usr.User{}
			u.Unmarshal(it.Item().KeyCopy(nil), nil)
			if filter(u) {
				users = append(users, *u)
			}
		}
		return nil
	})
	return
}

func (udb *UDB) IsExisting(uname string, onlyActive bool) bool {
	if onlyActive {
		_, okActive, err := udb.LoadUser(uname, true)
		lk.WarnOnErr("%v", err)
		return okActive
	}
	_, okActive, err := udb.LoadUser(uname, true)
	lk.WarnOnErr("%v", err)
	_, okDorm, err := udb.LoadUser(uname, false)
	lk.WarnOnErr("%v", err)
	return okActive || okDorm
}
