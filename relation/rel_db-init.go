package relation

import (
	"path/filepath"
	"sync"

	"github.com/dgraph-io/badger/v4"
	lk "github.com/digisan/logkit"
)

type DBGrp struct {
	sync.Mutex
	Rel *badger.DB
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
				Rel: open(filepath.Join(dir, "relation")),
			}
		})
	}
	return DbGrp
}

func CloseDB() {
	DbGrp.Lock()
	defer DbGrp.Unlock()

	if DbGrp.Rel != nil {
		lk.FailOnErr("%v", DbGrp.Rel.Close())
		DbGrp.Rel = nil
	}
}
