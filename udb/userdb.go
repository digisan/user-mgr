package udb

var (
	UserDB *UDB // global, for using
)

func OpenSession(udbPath string) {
	if UserDB == nil {
		UserDB = GetDB(udbPath)
	}
}

func CloseSession() {
	if UserDB != nil {
		UserDB.Close()
		UserDB = nil
	}
}
