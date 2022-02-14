package main

import (
	"fmt"

	lk "github.com/digisan/logkit"
	si "github.com/digisan/user-mgr/sign-in"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
)

func main() {

	udb.OpenUserStorage("../../data/user")
	defer udb.CloseUserStorage()

	// get [user] from GET

	// Will be GET header
	user := &usr.User{
		Active:     "",
		UName:      "QMiao",
		Email:      "",
		Name:       "",
		Password:   "pa55a@aD20TTTTT",
		Regtime:    "",
		Phone:      "",
		Addr:       "",
		SysRole:    "",
		MemLevel:   "",
		MemExpire:  "",
		NationalID: "",
		Gender:     "",
		Position:   "",
		Title:      "",
		Employer:   "",
		Avatar:     "",
	}

	lk.FailOnErr("%v", si.UserExists(user))
	lk.FailOnErrWhen(!si.PwdOK(user), "%v", fmt.Errorf("incorrect password"))
	lk.FailOnErr("%v", si.Trail(user.UName))

	fmt.Println("Login OK")
}
