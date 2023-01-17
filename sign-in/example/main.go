package main

import (
	"fmt"
	"time"

	lk "github.com/digisan/logkit"
	si "github.com/digisan/user-mgr/sign-in"
	u "github.com/digisan/user-mgr/user"
)

func main() {

	u.InitDB("../../data/user")
	defer u.CloseDB()

	// get [user] from GET

	// Will be GET header
	user := &u.User{
		u.Core{
			UName:    "Qing Miao",
			Email:    "4987346@qq.com",
			Password: "*pa55a@aD20TTTTT",
		},
		u.Profile{
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
		u.Admin{
			Regtime:   time.Now().Truncate(time.Second),
			Active:    true,
			Certified: false,
			Official:  false,
			SysRole:   "admin",
			MemLevel:  1,
			MemExpire: time.Time{},
			Tags:      "",
		},
	}

	lk.FailOnErr("%v", si.UserStatusIssue(user))
	lk.FailOnErrWhen(!si.PwdOK(user), "%v", fmt.Errorf("incorrect password"))
	lk.FailOnErr("%v", si.Trail(user.UName))

	fmt.Println("Login OK")
}
