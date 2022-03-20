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

// cache for fast fetching
var tmpUserPool = &sync.Map{}

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

func (db *UDB) GetOnline(uname string) (time.Time, error) {
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

func (db *UDB) RefreshOnline(uname string) error {
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

func (db *UDB) RmOnline(uname string) (err error) {
	// remove users cache when removing online user
	defer func() {
		if err == nil {
			if u, ok, err := db.LoadAnyUser(uname); err == nil && ok {
				tmpUserPool.Delete(u.UName)
				tmpUserPool.Delete(u.Email)
				tmpUserPool.Delete(u.Phone)
			}
		}
	}()

	db.Lock()
	defer db.Unlock()

	// we need err in defer()
	err = db.dbOnline.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(uname))
	})
	return
}

func (db *UDB) OnlineUsers() (unames []string, err error) {
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

func (db *UDB) UpdateUser(user *usr.User) (err error) {
	// update cache
	defer func() {
		if err == nil {
			tmpUserPool.Store(user.UName, user)
			tmpUserPool.Store(user.Email, user)
			tmpUserPool.Store(user.Phone, user)
		}
	}()

	db.Lock()
	defer db.Unlock()

	// remove all existing items
	if err = db.RemoveUser(user.UName, false, false); err != nil {
		return err
	}
	err = db.dbReg.Update(func(txn *badger.Txn) error {
		if forKey, forValue := user.Marshal(); len(forKey) > 0 || len(forValue) > 0 {
			return txn.Set(forKey, forValue)
		}
		return nil
	})
	return err
}

func (db *UDB) LoadUser(uname string, active bool) (*usr.User, bool, error) {

	// cache fetch & update
	if user, ok := tmpUserPool.Load(uname); ok {
		if u := user.(*usr.User); u.Email != "" {
			if (active && u.IsActive()) || (!active && !u.IsActive()) {
				return u, ok, nil
			}
		}
	}

	u := &usr.User{}
	var err error

	defer func() {
		if err == nil && u.Email != "" {
			tmpUserPool.Store(u.UName, u)
			tmpUserPool.Store(u.Email, u)
			tmpUserPool.Store(u.Phone, u)
		}
	}()

	///////////////////////////////////////////////////

	db.Lock()
	defer db.Unlock()

	prefix := []byte("T" + usr.SEP + uname + usr.SEP)
	if !active {
		prefix = []byte("F" + usr.SEP + uname + usr.SEP)
	}

	err = db.dbReg.View(func(txn *badger.Txn) error {
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
	return u, err == nil && u.Email != "", err
}

func (db *UDB) LoadActiveUser(uname string) (*usr.User, bool, error) {
	return db.LoadUser(uname, true)
}

func (db *UDB) LoadAnyUser(uname string) (*usr.User, bool, error) {
	uA, okA, errA := db.LoadUser(uname, true)
	uD, okD, errD := db.LoadUser(uname, false)
	var u *usr.User
	if okA {
		u = uA
	} else if okD {
		u = uD
	}
	var err error
	if errA != nil {
		err = errA
	} else if errD != nil {
		err = errD
	}
	return u, err == nil && (okA || okD), err
}

func (db *UDB) LoadUserByUniProp(propName, propVal string, active bool) (*usr.User, bool, error) {

	// cache fetch & update
	if user, ok := tmpUserPool.Load(propVal); ok {
		if u := user.(*usr.User); u.Email != "" {
			if (active && u.IsActive()) || (!active && !u.IsActive()) {
				return u, ok, nil
			}
		}
	}

	u := &usr.User{}
	var err error

	defer func() {
		if err == nil && u.Email != "" {
			tmpUserPool.Store(u.UName, u)
			tmpUserPool.Store(u.Email, u)
			tmpUserPool.Store(u.Phone, u)
		}
	}()

	///////////////////////////////////////////////////

	users, err := db.ListUsers(func(u *usr.User) bool {
		flag := u.IsActive()
		if !active {
			flag = !u.IsActive()
		}
		switch propName {
		case "uname", "Uname":
			return flag && u.UName == propVal
		case "email", "Email":
			return flag && u.Email == propVal
		case "phone", "Phone":
			return flag && u.Phone == propVal
		default:
			return false
		}
	})
	if len(users) > 0 {
		u = users[0]
		return u, err == nil && u.Email != "", err
	}
	return u, false, err
}

func (db *UDB) LoadActiveUserByUniProp(propName, propVal string) (*usr.User, bool, error) {
	return db.LoadUserByUniProp(propName, propVal, true)
}

func (db *UDB) LoadAnyUserByUniProp(propName, propVal string) (*usr.User, bool, error) {
	uA, okA, errA := db.LoadUserByUniProp(propName, propVal, true)
	uD, okD, errD := db.LoadUserByUniProp(propName, propVal, false)
	var u *usr.User
	if okA {
		u = uA
	} else if okD {
		u = uD
	}
	var err error
	if errA != nil {
		err = errA
	} else if errD != nil {
		err = errD
	}
	return u, err == nil && (okA || okD), err
}

func (db *UDB) RemoveUser(uname string, lock, rmCache bool) error {
	if rmCache {
		defer func() {
			if u, ok, err := db.LoadAnyUser(uname); err == nil && ok {
				tmpUserPool.Delete(u.UName)
				tmpUserPool.Delete(u.Email)
				tmpUserPool.Delete(u.Phone)
			}
		}()
	}

	if lock {
		db.Lock()
		defer db.Unlock()
	}

	prefixList := [][]byte{
		[]byte("T" + usr.SEP + uname + usr.SEP),
		[]byte("F" + usr.SEP + uname + usr.SEP),
	}
	return db.dbReg.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for _, prefix := range prefixList {
			if it.Seek(prefix); it.ValidForPrefix(prefix) {
				if err := txn.Delete(it.Item().KeyCopy(nil)); err != nil {
					return err
				}
			}
		}
		return nil
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

func (db *UDB) UserExists(uname, email string, activeOnly bool) bool {
	if activeOnly {
		// check uname
		_, ok, err := db.LoadUser(uname, true)
		lk.WarnOnErr("%v", err)
		if ok {
			return ok
		}
		// check email
		_, ok, err = db.LoadActiveUserByUniProp("email", email)
		lk.WarnOnErr("%v", err)
		return ok

	} else {
		// check uname
		_, ok, err := db.LoadAnyUser(uname)
		lk.WarnOnErr("%v", err)
		if ok {
			return ok
		}
		// check email
		_, ok, err = db.LoadAnyUserByUniProp("email", email)
		lk.WarnOnErr("%v", err)
		return ok
	}
}

func (db *UDB) ActivateUser(uname string, flag bool) (*usr.User, bool, error) {
	return db.SetUserBoolField(uname, "active", flag)
}

func (db *UDB) OfficializeUser(uname string, flag bool) (*usr.User, bool, error) {
	return db.SetUserBoolField(uname, "official", flag)
}

func (db *UDB) SetUserBoolField(uname, field string, flag bool) (*usr.User, bool, error) {
	val := strings.ToUpper(fmt.Sprint(flag))[:1]
	u, ok, err := db.LoadAnyUser(uname)
	if err == nil {
		if ok {
			switch field {
			case "Active", "active", "ACTIVE":
				u.Active = val
			case "Official", "official", "OFFICIAL":
				u.Official = val
			default:
				lk.FailOnErr("%v", fmt.Errorf("[%s] is unsupported setting BoolField", field))
			}
			err = db.UpdateUser(u)
			return u, err == nil, err
		}
		return nil, false, fmt.Errorf("couldn't find [%s] for setting [%s]", uname, field)
	}
	return nil, false, err
}
