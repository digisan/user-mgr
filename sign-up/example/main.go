package main

import (
	"fmt"
	"time"

	lk "github.com/digisan/logkit"
	. "github.com/digisan/user-mgr/db"
	su "github.com/digisan/user-mgr/sign-up"
	. "github.com/digisan/user-mgr/user"
	ur "github.com/digisan/user-mgr/user/registered"
)

func main() {

	lk.WarnDetail(false)

	InitDB("../../server-example/data/user")
	defer CloseDB()

	// get [user] from POST

	// Will be POST header
	user := &ur.User{
		Core: ur.Core{
			UName:    "Qing.Miao",
			Email:    "4987346@qq.com",
			Password: "*pa55a@aD20TTTTT",
		},
		Profile: ur.Profile{
			Name:           "A boy has no name",
			Phone:          "111111111",
			Country:        "",
			City:           "",
			Addr:           "",
			PersonalIDType: "",
			PersonalID:     "",
			Gender:         "",
			DOB:            "",
			Position:       "",
			Title:          "",
			Employer:       "ABCDEFG",
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

	su.SetValidator(map[string]func(o any, v any) ResultOk{
		ur.Employer: func(o, v any) ResultOk {
			ok := len(v.(string)) > 6
			return NewResultOk(ok, "at least 6 length for employer")
		},
	})

	if err := su.CheckInput(user, ur.Phone); err != nil {
		lk.WarnOnErr("%v", err)
		return
	}

	if err := su.CheckEmail(user); err != nil {
		lk.WarnOnErr("%v", err)
		return
	}

	fmt.Println("Input verification code in your email:", user.Email)
	incode := ""
	fmt.Scanf("%s", &incode)

	// get [incode] from POST
	if _, err := su.VerifyCode(user.UName, incode); err != nil {
		lk.Warn("Sign-Up failed: ", err)
		return
	}

	// double check input before storing
	if err := su.CheckInput(user); err != nil {
		lk.WarnOnErr("%v", err)
		return
	}

	// store into db
	lk.FailOnErr("%v", su.Store(user))

	lk.Log("Sign-Up OK")
}
