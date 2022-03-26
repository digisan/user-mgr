package user

import (
	"fmt"
	"testing"
	"time"
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
			MemLevel:  "",
			MemExpire: time.Time{},
			Regtime:   time.Now(),
			Official:  false,
			Tags:      "",
		},
	}
	fmt.Println(user)

	claims := MakeUserClaims(user)
	fmt.Println(claims.GenToken())
	fmt.Println(claims.GenToken())

	//////////////////////////////

	user = &User{
		Core{
			UName:    "unique-user-name",
			Email:    "hello@abc.com",
			Password: "123456789a",
			key:      [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		Profile{
			Name:       "test-name",
			Phone:      "",
			Country:    "",
			City:       "",
			Addr:       "",
			PersonalID: "9876543210",
			Gender:     "",
			DOB:        "",
			Position:   "professor",
			Title:      "",
			Employer:   "",
			Bio:        "",
			AvatarType: "image/png",
			Avatar:     []byte("******"),
		},
		Admin{
			Active:    true,
			SysRole:   "",
			MemLevel:  "",
			MemExpire: time.Time{},
			Regtime:   time.Now(),
			Official:  false,
			Tags:      "",
		},
	}
	fmt.Println(user)

	claims = MakeUserClaims(user)
	fmt.Println(claims.GenToken())
	fmt.Println(claims.GenToken())

	fmt.Println()

	mUserToken.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})

	fmt.Println()
}
