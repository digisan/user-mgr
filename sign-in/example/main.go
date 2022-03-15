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
		UName:      "Qing Miao",
		Email:      "",
		Name:       "",
		Password:   "*pa55a@aD20TTTTT",
		Regtime:    "",
		Official:   "",
		Phone:      "",
		Country:    "",
		City:       "",
		Addr:       "",
		SysRole:    "",
		MemLevel:   "",
		MemExpire:  "",
		NationalID: "",
		Gender:     "",
		DOB:        "",
		Position:   "",
		Title:      "",
		Employer:   "",
		Bio:        "",
		Tags:       "",
		AvatarType: "",
		Avatar:     []byte(""),
	}

	lk.FailOnErr("%v", si.CheckUserExists(user))
	lk.FailOnErrWhen(!si.PwdOK(user), "%v", fmt.Errorf("incorrect password"))
	lk.FailOnErr("%v", si.Trail(user.UName))

	fmt.Println("Login OK")
}
