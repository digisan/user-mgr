package user

import (
	"fmt"
	"testing"
)

func TestClaims(t *testing.T) {

	user := &User{
		Active:     "T",
		UName:      "unique-user-name",
		Email:      "hello@abc.com",
		Name:       "test-name",
		Password:   "123456789a",
		Regtime:    "",
		Phone:      "",
		Addr:       "",
		SysRole:    "",
		MemLevel:   "",
		MemExpire:  "",
		NationalID: "",
		Gender:     "",
		Position:   "",
		Title:      "",
		Employer:   "",
		Avatar:     "",
		key:        "",
	}
	fmt.Println(user)

	claims := MakeUserClaims(user)
	fmt.Println(claims.GenToken())
	fmt.Println(claims.GenToken())

	//////////////////////////////

	user = &User{
		Active:     "T",
		UName:      "unique-user-name-1",
		Email:      "hello@abc.com",
		Name:       "test-name-1",
		Password:   "123456789a",
		Regtime:    "",
		Phone:      "",
		Addr:       "",
		SysRole:    "",
		MemLevel:   "",
		MemExpire:  "",
		NationalID: "",
		Gender:     "",
		Position:   "",
		Title:      "",
		Employer:   "",
		Avatar:     "",
		key:        "",
	}
	fmt.Println(user)

	claims = MakeUserClaims(user)
	fmt.Println(claims.GenToken())
	fmt.Println(claims.GenToken())

	fmt.Println()

	mUserToken.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})

	fmt.Println()
}
