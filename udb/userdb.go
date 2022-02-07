package udb

var (
	UserDB *UDB // global, for using
)

// initiate [udb.UserDB] for using
func OpenUserStorage(udbPath string) {
	if UserDB == nil {
		UserDB = getDB(udbPath)
	}
}

func CloseUserStorage() {
	if UserDB != nil {
		UserDB.close()
		UserDB = nil
	}
}
