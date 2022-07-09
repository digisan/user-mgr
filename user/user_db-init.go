package user

import (
	"path/filepath"
	"sync"

	"github.com/dgraph-io/badger/v3"
	lk "github.com/digisan/logkit"
)

type DBGrp struct {
	sync.Mutex
	Reg    *badger.DB
	Online *badger.DB
}

var (
	onceDB sync.Once // do once
	DbGrp  *DBGrp    // global, for keeping single instance
)

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
	if DbGrp == nil {
		onceDB.Do(func() {
			DbGrp = &DBGrp{
				Reg:    open(filepath.Join(dir, "registration")),
				Online: open(filepath.Join(dir, "online")),
			}
		})
	}
	return DbGrp
}

func CloseDB() {
	DbGrp.Lock()
	defer DbGrp.Unlock()

	if DbGrp.Reg != nil {
		lk.FailOnErr("%v", DbGrp.Reg.Close())
		DbGrp.Reg = nil
	}
	if DbGrp.Online != nil {
		lk.FailOnErr("%v", DbGrp.Online.Close())
		DbGrp.Online = nil
	}
}
