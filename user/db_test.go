package user

import (
	"fmt"
	"testing"
	"time"

	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/db"
	. "github.com/digisan/user-mgr/user/registered"
)

const dbPath = "../data/db/user"

func TestListUser(t *testing.T) {

	db.InitDB(dbPath)
	defer db.CloseDB()

	users, err := ListUser(func(u *User) bool {
		return u.IsActive() || !u.IsActive()
	})
	lk.FailOnErr("%v", err)

	for _, user := range users {
		fmt.Println(user)
		fmt.Printf("%v since registered \n", user.SinceJoined())
	}
}

func TestSaveUser(t *testing.T) {

	db.InitDB(dbPath)
	defer db.CloseDB()

	u := &User{
		Core: Core{
			UName:    "unique-user-name",
			Email:    "hello@abc.com",
			Password: "123456789a",
		},
		Profile: Profile{
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
		Admin: Admin{
			RegTime:   time.Now().Truncate(time.Second),
			Active:    false,
			Certified: false,
			Official:  false,
			SysRole:   "",
			MemLevel:  0,
			MemExpire: time.Time{},
			Notes:     "",
			Status:    "",
		},
	}

	lk.FailOnErr("%v", UpdateUser(u))
}

func TestOpen(t *testing.T) {

	db.InitDB(dbPath)
	defer db.CloseDB()

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

	db.InitDB(dbPath)
	defer db.CloseDB()

	lk.FailOnErr("%v", RemoveUser("unique-user-name", true))

	u1, ok, err := LoadActiveUserByUniProp("email", "hello@abc.com")
	fmt.Println(u1, ok, err)

	fmt.Println("---------------------------")

	u2, ok, err := LoadActiveUserByUniProp("phone", "111111111")
	fmt.Println(u2, ok, err)
}

func TestExisting(t *testing.T) {

	db.InitDB(dbPath)
	defer db.CloseDB()

	fmt.Println("---", UserExists("unique-user-name", "", false))
	fmt.Println("---", UserExists("", "hello@abc.com", false))
}

///////////////////////////////////////////////////////////////////

func TestOnlines(t *testing.T) {

	db.InitDB(dbPath)
	defer db.CloseDB()

	users, err := OnlineUsers()
	lk.FailOnErr("%v", err)

	for _, u := range users {
		fmt.Println(u)
	}
}

func TestRefreshOnline(t *testing.T) {

	db.InitDB(dbPath)
	defer db.CloseDB()

	fmt.Println(RefreshOnline("uname1"))
}

func TestGetOnline(t *testing.T) {

	db.InitDB(dbPath)
	defer db.CloseDB()

	fmt.Println(GetOnline("uname1"))
}

func TestRmOnline(t *testing.T) {

	db.InitDB(dbPath)
	defer db.CloseDB()

	fmt.Println(RmOnline("uname1"))
}

func TestUpdateOnlineUser(t *testing.T) {

	db.InitDB(dbPath)
	defer db.CloseDB()

	RefreshOnline("a")
	time.Sleep(1 * time.Second)
	RefreshOnline("b")
	time.Sleep(1 * time.Second)
	RefreshOnline("c")
	time.Sleep(1 * time.Second)

	users, err := OnlineUsers()
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	u, _ := GetOnline("a")
	fmt.Println(*u)

	time.Sleep(3 * time.Second)

	if time.Since(u.Tm) > 2*time.Second {
		fmt.Println("more than 2 seconds")
	} else {
		fmt.Println("less than 2 seconds")
	}
}
