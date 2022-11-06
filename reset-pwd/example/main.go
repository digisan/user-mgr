package main

import (
	"fmt"
	"time"

	lk "github.com/digisan/logkit"
	rp "github.com/digisan/user-mgr/reset-pwd"
	su "github.com/digisan/user-mgr/sign-up"
	u "github.com/digisan/user-mgr/user"
)

func main() {

	u.InitDB("../../data/user")
	defer u.CloseDB()

	// get [user] from GET

	// Will be GET header
	usr := u.User{
		u.Core{
			UName:    "QMiao",
			Email:    "cdutwhu@outlook.com",
			Password: "",
		},
		u.Profile{
			Name:           "",
			Phone:          "",
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
			Avatar:         []byte{},
		},
		u.Admin{
			Regtime:   time.Now().Truncate(time.Second),
			Active:    true,
			Certified: false,
			Official:  false,
			SysRole:   "",
			MemLevel:  0,
			MemExpire: time.Time{},
			Tags:      "",
		},
	}

	if err := rp.CheckUserExists(&usr); err != nil {
		lk.Log("%v", err)
		return
	}
	if !rp.EmailOK(&usr) {
		lk.Log("%s's email [%s] is different from your sign-up one", usr.UName, usr.Email)
		return
	}

	if err := su.ChkEmail(&usr); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Input verification code in you email:", usr.Email)
	incode := ""
	fmt.Scanf("%s", &incode)
	// get [incode] from POST

	user, err := su.VerifyCode(usr.UName, incode)
	if err != nil {
		fmt.Println("Email verification failed:", err)
		return
	}

	// user, _, err = u.LoadActiveUser(user.UName)
	// lk.FailOnErr("%v", err)
	/////

AGAIN:
	fmt.Println("Input new password")
	pwdUpdated := ""
	fmt.Scanf("%s", &pwdUpdated)
	// get [pwdUpdated] from POST

	if rst := su.ChkPwd(pwdUpdated); rst.OK {
		user.Password = pwdUpdated
	} else {
		fmt.Println("invalid new password")
		goto AGAIN
	}

	// store into db
	if err := su.Store(&usr); err != nil {
		fmt.Println(err)
	}

	fmt.Println("OK")
}
