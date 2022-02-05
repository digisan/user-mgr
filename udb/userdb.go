package udb

var (
	UserDB *UDB // global, for using
)

// init udb.UserDB
func OpenSession(udbPath string) {
	if UserDB == nil {
		UserDB = getDB(udbPath)
	}
}

func CloseSession() {
	if UserDB != nil {
		UserDB.Close()
		UserDB = nil
	}
}
