package main

import (
	"fmt"

	lk "github.com/digisan/logkit"
	su "github.com/digisan/user-mgr/sign-up"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	vf "github.com/digisan/user-mgr/user/valfield"
)

func main() {

	lk.WarnDetail(false)

	udb.OpenUserStorage("../../data/user")
	defer udb.CloseUserStorage()

	// get [user] from POST

	// Will be POST header
	user := &usr.User{
		Active:     "T",
		UName:      "Qing Miao",
		Email:      "4987346@qq.com",
		Name:       "A boy has no name",
		Password:   "pa55a@aD20TTTTT",
		Regtime:    "",
		Official:   "F",
		Phone:      "1",
		Country:    "",
		City:       "",
		Addr:       "",
		SysRole:    "admin",
		MemLevel:   "1",
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
		Avatar:     []byte("abcdefg**********"),
	}

	su.SetValidator(nil)

	if err := su.ChkInput(user, vf.Phone); err != nil { // vf.Phone
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
	if _, err := su.VerifyCode(user.UName, incode); err != nil {
		fmt.Println("Sign-Up failed:", err)
		return
	}

	// store into db
	if err := su.Store(user); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Sign-up OK")
}
