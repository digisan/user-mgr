package udb

import (
	"fmt"
	"testing"
	"time"

	lk "github.com/digisan/logkit"
	usr "github.com/digisan/user-mgr/user"
)

const dbPath = "../data/user"

func TestOpen(t *testing.T) {
	udb := GetDB(dbPath)
	defer udb.Close()

	u := usr.User{
		Active:   "T",
		UName:    "unique-name",
		Email:    "hello@abc.net",
		Name:     "test-name",
		Password: "this is my password",
		Phone:    "123",
	}
	lk.FailOnErr("%v", udb.UpdateUser(u))

	u.Activate(false)
	lk.FailOnErr("%v", udb.UpdateUser(u))
}

func TestRemove(t *testing.T) {
	udb := GetDB(dbPath)
	defer udb.Close()

	lk.FailOnErr("%v", udb.RemoveUser("unique-name"))
}

func TestLoad(t *testing.T) {
	udb := GetDB(dbPath)
	defer udb.Close()

	user, ok, err := udb.LoadUser("unique-name", false)
	lk.FailOnErr("%v", err)

	fmt.Println(ok)
	fmt.Println(user)

	// check sign in
	// user.Password == "abc"
}

func TestListUsers(t *testing.T) {
	udb := GetDB(dbPath)
	defer udb.Close()

	users, err := udb.ListUsers(func(u *usr.User) bool {
		return u.IsActive() || !u.IsActive()
	})
	lk.FailOnErr("%v", err)

	for _, u := range users {
		fmt.Println(u)
		t := &time.Time{}
		t.UnmarshalText([]byte(u.Regtime))
		fmt.Println("Been regstered for", time.Since(*t))
	}
}

func TestExisting(t *testing.T) {
	udb := GetDB(dbPath)
	defer udb.Close()

	fmt.Println(udb.IsExisting("unique-name", false))
}

///////////////////////////////////////////////////////////////

func TestUpdateOnlineUser(t *testing.T) {
	udb := GetDB(dbPath)
	defer udb.Close()

	udb.UpdateOnlineUser("a")
	udb.UpdateOnlineUser("b")
	udb.UpdateOnlineUser("c")

	users, err := udb.ListOnlineUsers()
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	tm, _ := udb.LoadOnlineUser("a")
	fmt.Println(tm)

	time.Sleep(3 * time.Second)

	if time.Since(tm) > 2*time.Second {
		fmt.Println("more than 2 seconds")
	} else {
		fmt.Println("less than 2 seconds")
	}
}
