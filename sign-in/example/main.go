package main

import (
	"fmt"
	"time"

	lk "github.com/digisan/logkit"
	. "github.com/digisan/user-mgr/cst"
	. "github.com/digisan/user-mgr/db"
	si "github.com/digisan/user-mgr/sign-in"
	ur "github.com/digisan/user-mgr/user/registered"
)

func main() {

	InitDB("../../server-example/data/user")
	defer CloseDB()

	// get [user] from GET

	// Will be GET header
	user := &ur.User{
		Core: ur.Core{
			UName:    "Qing Miao",
			Email:    "4987346@qq.com",
			Password: "*pa55a@aD20TTTTT",
		},
		Profile: ur.Profile{
			Name:           "A boy has no name",
			Phone:          "11",
			Country:        "",
			City:           "",
			Addr:           "",
			PersonalIDType: "",
			PersonalID:     "",
			Gender:         "",
			DOB:            "",
			Position:       "",
			Title:          "",
			Employer:       "",
			Bio:            "",
			AvatarType:     "",
			Avatar:         []byte("abcdefg**********"),
		},
		Admin: ur.Admin{
			RegTime:   time.Now().Truncate(time.Second),
			Active:    true,
			Certified: false,
			Official:  false,
			SysRole:   "admin",
			MemLevel:  1,
			MemExpire: time.Time{},
			Notes:     "",
			Status:    "",
		},
	}

	lk.FailOnErr("%v", si.UserStatusIssue(user))
	lk.FailOnErrWhen(!si.PwdOK(user), "%v", Err(ERR_USER_PWD_INCORRECT))
	lk.FailOnErr("%v", si.Hail(user.UName))

	fmt.Println("Login OK")
}
