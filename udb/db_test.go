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
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	u := usr.User{
		Active:   "T",
		UName:    "unique-name",
		Email:    "hello@abc.net",
		Name:     "test-name",
		Password: "this is my password",
		Phone:    "123",
	}
	lk.FailOnErr("%v", UserDB.UpdateUser(u))

	done, err := UserDB.ActivateUser("unique-name", false)
	lk.FailOnErr("%v - %v", done, err)
}

func TestRemove(t *testing.T) {
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	lk.FailOnErr("%v", UserDB.RemoveUser("unique-name"))
}

func TestLoad(t *testing.T) {
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	user, ok, err := UserDB.LoadUser("unique-name", false)
	lk.FailOnErr("%v", err)

	fmt.Println(ok)
	fmt.Println(user)

	// check sign in
	// user.Password == "abc"
}

func TestListUsers(t *testing.T) {
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	users, err := UserDB.ListUsers(func(u *usr.User) bool {
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
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	fmt.Println(UserDB.IsExisting("unique-name", false))
}

///////////////////////////////////////////////////////////////

func TestUpdateOnlineUser(t *testing.T) {
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	UserDB.UpdateOnlineUser("a")
	UserDB.UpdateOnlineUser("b")
	UserDB.UpdateOnlineUser("c")

	users, err := UserDB.ListOnlineUsers()
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	tm, _ := UserDB.LoadOnlineUser("a")
	fmt.Println(tm)

	time.Sleep(3 * time.Second)

	if time.Since(tm) > 2*time.Second {
		fmt.Println("more than 2 seconds")
	} else {
		fmt.Println("less than 2 seconds")
	}
}
