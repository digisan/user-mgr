package relation

import (
	"fmt"
	"testing"
	"time"

	u "github.com/digisan/user-mgr/user"
)

func TestUserInit(t *testing.T) {
	u.InitDB("./data")
	defer u.CloseDB()

	usr := u.User{
		Core: u.Core{
			UName:    "qing",
			Email:    "",
			Password: "",
		},
		Profile: u.Profile{
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
		Admin: u.Admin{
			Regtime:   time.Time{},
			Active:    false,
			Certified: false,
			Official:  false,
			SysRole:   "",
			MemLevel:  0,
			MemExpire: time.Time{},
			Tags:      "",
		},
	}
	u.UpdateUser(&usr)
}

func TestUserCheck(t *testing.T) {
	u.InitDB("./data")
	defer u.CloseDB()

	fmt.Println(u.ListUser(nil))
}

func TestListRel(t *testing.T) {
	InitDB("./data")
	defer CloseDB()

	fmt.Println("followers:", ListRel("qing", FOLLOWER))
}

func TestRelAction(t *testing.T) {

	u.InitDB("./data")
	defer u.CloseDB()

	InitDB("./data")
	defer CloseDB()

	fmt.Println(RelAction("qing", DO_FOLLOW, "musk"))

	// RelAction("qing", DO_FOLLOW, "trump")
	// RelAction("qing", DO_UNFOLLOW, "musk")
	// RelAction("qing", DO_FOLLOW, "musk")

	// content := RelContent("qing", FOLLOWING)
	// fmt.Println(content)

	// rel := RelMgr("qing")
	// fmt.Println(rel.HasFollowing("musk"))

	// fmt.Println("-----------------------")

	// fmt.Println(RelMgr("musk"))

	// fmt.Println("-----------------------")

	// fmt.Println(RelMgr("trump"))

	// fmt.Println(" unfollow -----------------------")

	// RelAction("qing", DO_UNFOLLOW, "musk")
	// // RelAction(DO_UNFOLLOW, "qing", "trump")

	// fmt.Println(RelMgr("qing"))
	// fmt.Println(RelMgr("musk"))

	// RelAction("musk", DO_FOLLOW, "trump")

	// // fmt.Println(RelContent("trump", FOLLOWER))

	// fmt.Println(RelMgr("trump"))

	// RelAction("qing", DO_FOLLOW, "biden")
	// fmt.Println(RelMgr("qing"))

	// RelAction("qing", DO_BLOCK, "biden")
	// fmt.Println(RelMgr("qing"))
}
