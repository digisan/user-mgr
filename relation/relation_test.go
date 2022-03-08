package relation

import (
	"fmt"
	"testing"

	. "github.com/digisan/user-mgr/relation/enum"
)

func TestRelAction(t *testing.T) {
	OpenRelStorage("./data")
	defer CloseRelStorage()

	RelAction(DO_FOLLOW, "qing", "musk")
	RelAction(DO_FOLLOW, "qing", "trump")
	RelAction(DO_UNFOLLOW, "qing", "musk")
	RelAction(DO_FOLLOW, "qing", "musk")

	fmt.Println(RelMgr(FOLLOWING, "qing"))

	fmt.Println("-----------------------")

	fmt.Println(RelMgr(FOLLOWER, "musk"))

	fmt.Println("-----------------------")

	fmt.Println(RelMgr(FOLLOWER, "trump"))

	fmt.Println(" unfollow -----------------------")

	RelAction(DO_UNFOLLOW, "qing", "musk")
	RelAction(DO_UNFOLLOW, "qing", "trump")

	fmt.Println(RelMgr(FOLLOWING, "qing"))
	fmt.Println(RelMgr(FOLLOWER, "musk"))
	fmt.Println(RelMgr(FOLLOWER, "trump"))

}
