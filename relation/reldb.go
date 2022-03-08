package relation

import (
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v3"
	lk "github.com/digisan/logkit"
	. "github.com/digisan/user-mgr/relation/enum"
)

var once sync.Once

type RDB struct {
	sync.Mutex
	dbRel *badger.DB
}

var rDB *RDB

func open(dir string) *badger.DB {
	opt := badger.DefaultOptions("").WithInMemory(true)
	if dir != "" {
		opt = badger.DefaultOptions(dir)
	}
	db, err := badger.Open(opt)
	lk.FailOnErr("%v", err)
	return db
}

func getDB(dir string) *RDB {
	if rDB == nil {
		once.Do(func() {
			rDB = &RDB{
				dbRel: open(dir),
			}
		})
	}
	return rDB
}

func (db *RDB) close() {
	db.Lock()
	defer db.Unlock()

	lk.FailOnErr("%v", db.dbRel.Close())
}

////////////////////////////////////////////////////////////////////////

var (
	RelDB *RDB // global, for using
)

// initiate [RelDB] for using
func OpenRelStorage(rdbPath string) {
	if RelDB == nil {
		RelDB = getDB(rdbPath)
	}
}

func CloseRelStorage() {
	if RelDB != nil {
		RelDB.close()
		RelDB = nil
	}
}

////////////////////////////////////////////////////////////////////////

func (db *RDB) RemoveRel(flag int, uname string, lock bool) error {
	if lock {
		db.Lock()
		defer db.Unlock()
	}

	mPrefix := map[int][]byte{
		FOLLOWING: []byte("FI" + SEP + uname),
		FOLLOWER:  []byte("FR" + SEP + uname),
		BLOCKED:   []byte("B" + SEP + uname),
		MUTED:     []byte("M" + SEP + uname),
	}

	prefix, ok := mPrefix[flag]
	lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("invalid flag"))

	return db.dbRel.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		if it.Seek(prefix); it.ValidForPrefix(prefix) {
			return txn.Delete(it.Item().KeyCopy(nil))
		}
		return nil
	})
}

func (db *RDB) UpdateRel(flag int, rel *Rel) (err error) {
	db.Lock()
	defer db.Unlock()

	if err = db.RemoveRel(flag, rel.uname, false); err != nil {
		return err
	}
	return db.dbRel.Update(func(txn *badger.Txn) error {
		if forKey, forValue := rel.MarshalTo(flag); len(forKey) > 0 && len(forValue) > 0 {
			return txn.Set(forKey, forValue)
		}
		return nil
	})
}

func (db *RDB) LoadRel(flag int, uname string) (*Rel, bool, error) {
	db.Lock()
	defer db.Unlock()

	mPrefix := map[int][]byte{
		FOLLOWING: []byte("FI" + SEP + uname),
		FOLLOWER:  []byte("FR" + SEP + uname),
		BLOCKED:   []byte("B" + SEP + uname),
		MUTED:     []byte("M" + SEP + uname),
	}

	prefix, ok := mPrefix[flag]
	lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("invalid flag, only accept [FOLLOWING FOLLOWER BLOCKED MUTED]"))

	r := &Rel{}
	err := db.dbRel.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		if it.Seek(prefix); it.ValidForPrefix(prefix) {
			item := it.Item()
			k := item.Key()
			return item.Value(func(v []byte) error {
				r.UnmarshalFrom(k, v)
				return nil
			})
		}
		return nil
	})
	return r, r.uname != "", err
}

////////////////////////////////////////////////////////////////////////

func RelMgr(flag int, uname string) *Rel {
	if RelDB != nil {
		if rel, ok, err := RelDB.LoadRel(flag, uname); err == nil && ok {
			return rel
		}
		return &Rel{uname: uname}
	}
	return nil
}
