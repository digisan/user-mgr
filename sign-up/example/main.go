package main

import (
	"fmt"
	"time"

	lk "github.com/digisan/logkit"
	su "github.com/digisan/user-mgr/sign-up"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
)

func main() {

	lk.WarnDetail(false)

	udb.OpenUserStorage("../../data/user")
	defer udb.CloseUserStorage()

	// get [user] from POST

	// Will be POST header
	user := usr.User{
		Active:   "T",
		UName:    "QMiao",
		Email:    "cdutwhu@outlook.com",
		Name:     "A girl has no name",
		Password: "pa55a@aD20TTTTT",
		Phone:    "123456789",
		Regtime:  time.Now().UTC().Format(time.RFC3339),
		SysRole:  "admin",
		SysLevel: "1",
		Avatar:   "abcdefg",
	}

	su.SetValidator()

	if err := su.ChkInput(user); err != nil {
		lk.WarnOnErr("%v", err)
		return
	}

	if err := su.ChkEmail(user); err != nil {
		lk.WarnOnErr("%v", err)
		return
	}

	fmt.Println("Input verification code in your email:", user.Email)
	incode := ""
	fmt.Scanf("%s", &incode)

	// get [incode] from POST
	if err := su.VerifyCode(user, incode); err != nil {
		fmt.Println("Sign-Up failed:", err)
		return
	}

	// store into db
	if err := su.Store(user); err != nil {
		fmt.Println(err)
	}

	fmt.Println("OK")
}
