package db

import (
	"path/filepath"
	"sync"

	"github.com/dgraph-io/badger/v4"
	fd "github.com/digisan/gotk/file-dir"
	lk "github.com/digisan/logkit"
)

type DatabaseGroup struct {
	sync.Mutex
	Registered *badger.DB
	Online     *badger.DB
	Signing    *badger.DB
}

var (
	once  sync.Once      // do once
	DbGrp *DatabaseGroup // global, for keeping single instance
)

func open(dir string) *badger.DB {
	var opt badger.Options
	if len(dir) != 0 {
		if fd.DirExists(dir) {
			lk.Log("opening dir for BadgerDB: '%s'", dir)
		} else {
			lk.Log("creating dir for BadgerDB: '%s'", dir)
		}
		opt = badger.DefaultOptions(dir)
		opt.Logger = nil
	} else {
		opt = badger.DefaultOptions("").WithInMemory(true)
		lk.Log("badger is in-memory mode")
	}
	db, err := badger.Open(opt)
	lk.FailOnErr("%v", err)
	return db
}

// init global 'dbGrp'
func InitDB(dir string) *DatabaseGroup {
	if DbGrp == nil {
		once.Do(func() {
			DbGrp = &DatabaseGroup{
				Registered: open(filepath.Join(dir, "registered")),
				Online:     open(filepath.Join(dir, "online")),
				Signing:    open(filepath.Join(dir, "signing")),
			}
		})
	}
	return DbGrp
}

func CloseDB() {
	DbGrp.Lock()
	defer DbGrp.Unlock()

	if DbGrp.Registered != nil {
		lk.FailOnErr("%v", DbGrp.Registered.Close())
		DbGrp.Registered = nil
	}
	if DbGrp.Online != nil {
		lk.FailOnErr("%v", DbGrp.Online.Close())
		DbGrp.Online = nil
	}
	if DbGrp.Signing != nil {
		lk.FailOnErr("%v", DbGrp.Signing.Close())
		DbGrp.Signing = nil
	}
}
