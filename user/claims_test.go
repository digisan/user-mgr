package user

import (
	"fmt"
	"os"
	"testing"
	"time"

	lk "github.com/digisan/logkit"
	ur "github.com/digisan/user-mgr/user/registered"
)

func TestClaims(t *testing.T) {

	user := &ur.User{
		Core: ur.Core{
			UName:    "unique-user-name",
			Email:    "hello@abc.com",
			Password: "123456789a",
		},
		Profile: ur.Profile{
			Name:           "test-name",
			Phone:          "",
			Country:        "",
			City:           "",
			Addr:           "",
			PersonalIDType: "",
			PersonalID:     "9876543210",
			Gender:         "",
			DOB:            "",
			Position:       "professor",
			Title:          "",
			Employer:       "",
			Bio:            "",
			AvatarType:     "image/png",
			Avatar:         []byte("******"),
		},
		Admin: ur.Admin{
			RegTime:   time.Time{},
			Active:    true,
			Certified: false,
			Official:  false,
			SysRole:   "",
			MemLevel:  0,
			MemExpire: time.Time{},
			Notes:     "",
			Status:    "",
		},
	}
	fmt.Println(user)

	prvKey, err := os.ReadFile("../../server-example/cert/id_rsa")
	lk.FailOnErr("%v", err)

	pubKey, err := os.ReadFile("../../server-example/cert/id_rsa.pub")
	lk.FailOnErr("%v", err)

	claims := MakeUserClaims(user)
	ts, err := claims.GenerateToken(prvKey)
	lk.FailOnErr("%v", err)

	fmt.Println(ts)

	smToken.Range(func(key, value any) bool {
		fmt.Println("\nkey ---", key)
		fmt.Println("\nval ---", value.(*TokenInfo).value)
		return true
	})

	fmt.Println("---------------------------------------")

	userValidate, err := ValidateToken(user, ts, pubKey)
	lk.FailOnErr("%v", err)

	fmt.Printf("%+v", userValidate)
}
