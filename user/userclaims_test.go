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
			MemLevel:  0,
			MemExpire: time.Time{},
			RegTime:   time.Now(),
			Official:  false,
			Tags:      "",
		},
	}
	fmt.Println(user)

	claims := MakeClaims(user)
	fmt.Println(GenerateToken(claims))

	smToken.Range(func(key, value any) bool {
		fmt.Println("key ---", key)
		fmt.Println("val ---", value)
		return true
	})
}
