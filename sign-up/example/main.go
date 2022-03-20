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
