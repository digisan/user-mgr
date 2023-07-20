package relation

import (
	"fmt"
	"testing"
	"time"

	lk "github.com/digisan/logkit"
	. "github.com/digisan/user-mgr/db"
	u "github.com/digisan/user-mgr/user"
	ur "github.com/digisan/user-mgr/user/registered"
)

func TestUserInit(t *testing.T) {
	InitDB("../server-example/data")
	defer CloseDB()

	for _, uname := range []string{"qing", "musk", "trump"} {
		usr := ur.User{
			Core: ur.Core{
				UName:    uname,
				Email:    "abc@email.com",
				Password: "",
			},
			Profile: ur.Profile{
				Name:           "",
				Phone:          "",
				Country:        "",
				City:           "",
				Addr:           "",
				PersonalIDType: "",
				PersonalID:     "",
				Gender:         "",
				DOB:            "",
				Position:       "",
				Title:          "",
				Employer:       "",
				Bio:            "",
				AvatarType:     "",
				Avatar:         []byte{},
			},
			Admin: ur.Admin{
				RegTime:   time.Time{},
				Active:    true,
				Certified: false,
				Official:  false,
				SysRole:   "",
				MemLevel:  0,
				MemExpire: time.Time{},
				Notes:     "",
				Status:    "",
			},
		}
		u.UpdateUser(&usr)
	}
}

func TestUserCheck(t *testing.T) {
	InitDB("../server-example/data")
	defer CloseDB()

	fmt.Println(u.ListUser(nil))
}

func TestListRel(t *testing.T) {
	InitDB("./data")
	defer CloseDB()

	uname := "qing"

	names, err := ListRel(uname, FOLLOWER, true)
	lk.FailOnErr("%v", err)
	fmt.Println("followers:", names)
	fmt.Println("---------------------------------")
	fmt.Println(RelMgr(uname))
}

func TestClearRel(t *testing.T) {

	InitDB("../server-example/data")
	defer CloseDB()

	InitDB("./data")
	defer CloseDB()

	fmt.Println(RelAction("qing", UNFOLLOW, "ALL"))
}

func TestRelAction(t *testing.T) {

	InitDB("../server-example/data")
	defer CloseDB()

	InitDB("./data")
	defer CloseDB()

	// fmt.Println(RelAction("qing", FOLLOW, "trump"))
	// fmt.Println(RelAction("musk", FOLLOW, "qing"))

	// fmt.Println(RelAction("qing", FOLLOW, "trump"))
	fmt.Println(RelAction("qing", UNFOLLOW, "musk"))
	fmt.Println(RelAction("qing", FOLLOW, "musk"))

	// content := RelContent("qing", FOLLOWING)
	// fmt.Println(content)

	r, err := RelMgr("qing")
	lk.FailOnErr("%v", err)
	fmt.Println(RelMgr("qing"))
	fmt.Println("qing following musk:", r.HasFollowing("musk"))
	fmt.Println("qing following trump:", r.HasFollowing("trump"))

	fmt.Println("-----------------------")
	fmt.Println(RelMgr("musk"))

	fmt.Println("-----------------------")
	fmt.Println(RelMgr("trump"))

	// fmt.Println(" unfollow -----------------------")

	// RelAction("qing", UNFOLLOW, "musk")

	// fmt.Println(RelMgr("qing"))
	// fmt.Println(RelMgr("musk"))

	// RelAction("musk", FOLLOW, "trump")

	// fmt.Println(RelMgr("trump"))

	// RelAction("qing", FOLLOW, "biden")
	// fmt.Println(RelMgr("qing"))

	// RelAction("qing", BLOCK, "biden")
	// fmt.Println(RelMgr("qing"))
}
