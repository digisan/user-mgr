package registered

import (
	"fmt"
	"os"
	"testing"
	"time"

	lk "github.com/digisan/logkit"
)

func TestClaims(t *testing.T) {

	user := &User{
		Core{
			UName:    "unique-user-name",
			Email:    "hello@abc.com",
			Password: "123456789a",
			key:      [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		Profile{
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
		Admin{
			Active:    true,
			SysRole:   "",
			MemLevel:  0,
			MemExpire: time.Time{},
			RegTime:   time.Now(),
			Official:  false,
			Tags:      "",
		},
	}
	fmt.Println(user)

	prvKey, err := os.ReadFile("../server-example/cert/id_rsa")
	lk.FailOnErr("%v", err)

	pubKey, err := os.ReadFile("../server-example/cert/id_rsa.pub")
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

	userValidate, err := user.ValidateToken(ts, pubKey)
	lk.FailOnErr("%v", err)

	fmt.Printf("%+v", userValidate)
}
