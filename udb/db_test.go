package udb

import (
	"fmt"
	"testing"
	"time"

	lk "github.com/digisan/logkit"
	usr "github.com/digisan/user-mgr/user"
)

const dbPath = "../data/user"

func TestLoad(t *testing.T) {
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	user0, ok, err := UserDB.LoadUser("unique-user-name", true)
	lk.FailOnErr("%v", err)
	fmt.Println(ok)

	fmt.Println(user0)

	user1, ok, err := UserDB.LoadUserByUniProp("uname", "unique-user-name", true)
	lk.FailOnErr("%v", err)
	fmt.Println(ok)

	fmt.Println("=====", user0 == user1)
}

func TestOpen(t *testing.T) {
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	u := &usr.User{
		usr.Core{
			UName:    "unique-user-name",
			Email:    "hello@abc.com",
			Password: "123456789a",
		},
		usr.Profile{
			Name:       "test-name",
			Phone:      "111111111",
			Country:    "",
			City:       "",
			Addr:       "",
			NationalID: "9876543210",
			Gender:     "",
			DOB:        "",
			Position:   "professor",
			Title:      "",
			Employer:   "",
			Bio:        "",
			AvatarType: "image/png",
			Avatar:     []byte("******"),
		},
		usr.Admin{
			Regtime:   "",
			Active:    "T",
			SysRole:   "",
			MemLevel:  "",
			MemExpire: "",
			Official:  "",
			Tags:      "",
		},
	}

	lk.FailOnErr("%v", UserDB.UpdateUser(u))

	fmt.Println("---------------------------")

	u, done, err := UserDB.ActivateUser("unique-user-name", true)
	lk.WarnOnErr("------: %v - %v", done, err)
	fmt.Println(u)

	fmt.Println()

	u1, ok, err := UserDB.LoadActiveUserByUniProp("email", "hello@abc.com")
	fmt.Println(u1, ok, err)

	fmt.Println()

	u2, ok, err := UserDB.LoadActiveUserByUniProp("phone", "111111111")
	fmt.Println(u2, ok, err)
}

func TestRemove(t *testing.T) {
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	lk.FailOnErr("%v", UserDB.RemoveUser("unique-user-name", true, true))

	u1, ok, err := UserDB.LoadActiveUserByUniProp("email", "hello@abc.com")
	fmt.Println(u1, ok, err)

	fmt.Println()

	u2, ok, err := UserDB.LoadActiveUserByUniProp("phone", "111111111")
	fmt.Println(u2, ok, err)
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
		fmt.Println("been regstered for", time.Since(*t))
	}
}

func TestExisting(t *testing.T) {
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	fmt.Println("---", UserDB.UserExists("unique-user-name", "", false))
	fmt.Println("---", UserDB.UserExists("", "hello@abc.com", false))
}

///////////////////////////////////////////////////////////////

func TestUpdateOnlineUser(t *testing.T) {
	OpenUserStorage(dbPath)
	defer CloseUserStorage()

	UserDB.RefreshOnline("a")
	UserDB.RefreshOnline("b")
	UserDB.RefreshOnline("c")

	users, err := UserDB.OnlineUsers()
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	tm, _ := UserDB.GetOnline("a")
	fmt.Println(tm)

	time.Sleep(3 * time.Second)

	if time.Since(tm) > 2*time.Second {
		fmt.Println("more than 2 seconds")
	} else {
		fmt.Println("less than 2 seconds")
	}
}
