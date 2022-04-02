package main

import (
	"fmt"
	"time"

	lk "github.com/digisan/logkit"
	rp "github.com/digisan/user-mgr/reset-pwd"
	su "github.com/digisan/user-mgr/sign-up"
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
			UName:    "QMiao",
			Email:    "cdutwhu@outlook.com",
			Password: "",
		},
		usr.Profile{
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
		usr.Admin{
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

	if err := rp.CheckUserExists(user); err != nil {
		lk.Log("%v", err)
		return
	}
	if !rp.EmailOK(user) {
		lk.Log("%s's email [%s] is different from your sign-up one", user.UName, user.Email)
		return
	}

	if err := su.ChkEmail(user); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Input verification code in you email:", user.Email)
	incode := ""
	fmt.Scanf("%s", &incode)
	// get [incode] from POST

	user, err := su.VerifyCode(user.UName, incode)
	if err != nil {
		fmt.Println("Email verification failed:", err)
		return
	}

	// user, _, err = udb.UserDB.LoadActiveUser(user.UName)
	// lk.FailOnErr("%v", err)
	/////

AGAIN:
	fmt.Println("Input new password")
	pwdUpdated := ""
	fmt.Scanf("%s", &pwdUpdated)
	// get [pwdUpdated] from POST

	if su.ChkPwd(pwdUpdated, su.PwdLen) {
		user.Password = pwdUpdated
	} else {
		fmt.Println("invalid new password")
		goto AGAIN
	}

	// store into db
	if err := su.Store(user); err != nil {
		fmt.Println(err)
	}

	fmt.Println("OK")
}
