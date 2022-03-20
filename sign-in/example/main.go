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
		usr.Core{
			UName:    "Qing Miao",
			Email:    "4987346@qq.com",
			Password: "*pa55a@aD20TTTTT",
			Key:      [16]byte{},
		},
		usr.Profile{
			Name:       "A boy has no name",
			Phone:      "11",
			Country:    "",
			City:       "",
			Addr:       "",
			NationalID: "",
			Gender:     "",
			DOB:        "",
			Position:   "",
			Title:      "",
			Employer:   "",
			Bio:        "",
			AvatarType: "",
			Avatar:     []byte("abcdefg**********"),
		},
		usr.Admin{
			Regtime:   "",
			Active:    "T",
			SysRole:   "admin",
			MemLevel:  "1",
			MemExpire: "",
			Official:  "F",
			Tags:      "",
		},
	}

	lk.FailOnErr("%v", si.CheckUserExists(user))
	lk.FailOnErrWhen(!si.PwdOK(user), "%v", fmt.Errorf("incorrect password"))
	lk.FailOnErr("%v", si.Trail(user.UName))

	fmt.Println("Login OK")
}
