package user

import (
	"fmt"
	"testing"
	"time"

	lk "github.com/digisan/logkit"
)

const dbPath = "../data/db/user"

func TestListUser(t *testing.T) {

	InitDB(dbPath)
	defer CloseDB()

	users, err := ListUser(func(u *User) bool {
		return u.IsActive() || !u.IsActive()
	})
	lk.FailOnErr("%v", err)

	for _, user := range users {
		fmt.Println(user)
		fmt.Printf("regstered for %v\n\n", user.SinceJoined())
	}
}

func TestSaveUser(t *testing.T) {
	InitDB(dbPath)
	defer CloseDB()

	u := &User{
		Core{
			UName:    "unique-user-name",
			Email:    "hello@abc.com",
			Password: "123456789a",
		},
		Profile{
			Name:           "test-name",
			Phone:          "111111111",
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
			Regtime:   time.Now().Truncate(time.Second),
			Active:    false,
			Certified: false,
			Official:  false,
			SysRole:   "",
			MemLevel:  0,
			MemExpire: time.Time{},
			Tags:      "",
		},
	}

	lk.FailOnErr("%v", UpdateUser(u))
}

func TestOpen(t *testing.T) {

	InitDB(dbPath)
	defer CloseDB()

	fmt.Println("---------------------------")

	_, done, err := ActivateUser("unique-user-name", true)
	lk.WarnOnErr("------: %v - %v", done, err)
	// fmt.Println(u)

	_, done, err = OfficializeUser("unique-user-name", true)
	lk.WarnOnErr("------: %v - %v", done, err)
	// // fmt.Println(u)

	_, done, err = CertifyUser("unique-user-name", true)
	lk.WarnOnErr("------: %v - %v", done, err)
	// fmt.Println(u0)

	fmt.Println("---------------------------")

	u1, ok, err := LoadActiveUserByUniProp("email", "hello@abc.com")
	fmt.Println(u1, ok, err)

	fmt.Println("---------------------------")

	u2, ok, err := LoadActiveUserByUniProp("phone", "111111111")
	fmt.Println(u2, ok, err)
}

func TestRemove(t *testing.T) {

	InitDB(dbPath)
	defer CloseDB()

	lk.FailOnErr("%v", RemoveUser("unique-user-name", true))

	u1, ok, err := LoadActiveUserByUniProp("email", "hello@abc.com")
	fmt.Println(u1, ok, err)

	fmt.Println("---------------------------")

	u2, ok, err := LoadActiveUserByUniProp("phone", "111111111")
	fmt.Println(u2, ok, err)
}

func TestExisting(t *testing.T) {

	InitDB(dbPath)
	defer CloseDB()

	fmt.Println("---", UserExists("unique-user-name", "", false))
	fmt.Println("---", UserExists("", "hello@abc.com", false))
}

///////////////////////////////////////////////////////////////////

func TestOnlines(t *testing.T) {

	InitDB(dbPath)
	defer CloseDB()

	users, err := OnlineUsers()
	lk.FailOnErr("%v", err)

	for _, u := range users {
		fmt.Println(u)
	}
}

func TestRefreshOnline(t *testing.T) {

	InitDB(dbPath)
	defer CloseDB()

	fmt.Println(RefreshOnline("uname1"))
}

func TestGetOnline(t *testing.T) {

	InitDB(dbPath)
	defer CloseDB()

	fmt.Println(GetOnline("uname1"))
}

func TestRmOnline(t *testing.T) {

	InitDB(dbPath)
	defer CloseDB()

	fmt.Println(RmOnline("uname1"))
}
