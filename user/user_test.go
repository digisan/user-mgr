package user

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUser(t *testing.T) {
	user := &User{
		Active:   "T",
		UName:    "unique-user-name",
		Email:    "hello@abc.com",
		Name:     "test-name",
		Password: "123456789ab",
	}
	fmt.Println(user)

	info, key := user.Marshal()
	fmt.Println(user.key)

	user1 := &User{}
	user1.Unmarshal(info, key)
	fmt.Println(user1)

	fmt.Println("user == user1 :", user == user1)
	fmt.Println("reflect.DeepEqual(*user, *user1) :", reflect.DeepEqual(*user, *user1))
}

func TestMisc(t *testing.T) {
	fmt.Print(true)
	fmt.Print(false)
}
