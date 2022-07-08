package user

import (
	"path/filepath"
	"sync"

	"github.com/dgraph-io/badger/v3"
	lk "github.com/digisan/logkit"
)

var (
	onceDB sync.Once // do once
	dbGrp  *DBGrp    // global, for keeping single instance
)

type DBGrp struct {
	sync.Mutex
	dbReg    *badger.DB
	dbOnline *badger.DB
}

func open(dir string) *badger.DB {
	opt := badger.DefaultOptions("").WithInMemory(true)
	if dir != "" {
		opt = badger.DefaultOptions(dir)
		opt.Logger = nil
	}
	db, err := badger.Open(opt)
	lk.FailOnErr("%v", err)
	return db
}

// init global 'dbGrp'
func InitDB(dir string) *DBGrp {
	if dbGrp == nil {
		onceDB.Do(func() {
			dbGrp = &DBGrp{
				dbReg:    open(filepath.Join(dir, "registration")),
				dbOnline: open(filepath.Join(dir, "online")),
			}
		})
	}
	return dbGrp
}

func CloseDB() {
	dbGrp.Lock()
	defer dbGrp.Unlock()

	if dbGrp.dbReg != nil {
		lk.FailOnErr("%v", dbGrp.dbReg.Close())
		dbGrp.dbReg = nil
	}
	if dbGrp.dbOnline != nil {
		lk.FailOnErr("%v", dbGrp.dbOnline.Close())
		dbGrp.dbOnline = nil
	}
}
