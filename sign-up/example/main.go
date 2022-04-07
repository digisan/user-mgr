package main

import (
	"fmt"
	"time"

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
		},
		usr.Profile{
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
		usr.Admin{
			Regtime:   time.Now().Truncate(time.Second),
			Active:    true,
			Certified: false,
			Official:  false,
			SysRole:   "admin",
			MemLevel:  1,
			MemExpire: time.Time{},
			Tags:      "tag",
		},
	}

	su.SetValidator(map[string]func(o any, v any) usr.ValRst{
		vf.Employer: func(o, v any) usr.ValRst {
			ok := len(v.(string)) > 6
			return usr.NewValRst(ok, "at least 6 length for employer")
		},
	})

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
		lk.Warn("Sign-Up failed: ", err)
		return
	}

	// double check input before storing
	if err := su.ChkInput(user); err != nil {
		lk.WarnOnErr("%v", err)
		return
	}

	// store into db
	lk.FailOnErr("%v", su.Store(user))

	lk.Log("Sign-up OK")
}
