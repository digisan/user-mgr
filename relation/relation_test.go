package relation

import (
	"fmt"
	"testing"

	. "github.com/digisan/user-mgr/relation/enum"
)

func TestRelAction(t *testing.T) {
	OpenRelStorage("./data")
	defer CloseRelStorage()

	RelAction("qing", DO_FOLLOW, "musk")

	RelAction("qing", DO_FOLLOW, "trump")
	RelAction("qing", DO_UNFOLLOW, "musk")
	RelAction("qing", DO_FOLLOW, "musk")

	// content := RelContent("qing", FOLLOWING)
	// fmt.Println(content)

	rel := RelMgr("qing")
	fmt.Println(rel.HasFollowing("musk"))

	fmt.Println("-----------------------")

	fmt.Println(RelMgr("musk"))

	fmt.Println("-----------------------")

	fmt.Println(RelMgr("trump"))

	fmt.Println(" unfollow -----------------------")

	RelAction("qing", DO_UNFOLLOW, "musk")
	// RelAction(DO_UNFOLLOW, "qing", "trump")

	fmt.Println(RelMgr("qing"))
	fmt.Println(RelMgr("musk"))

	RelAction("musk", DO_FOLLOW, "trump")

	// fmt.Println(RelContent("trump", FOLLOWER))

	fmt.Println(RelMgr("trump"))

	RelAction("qing", DO_FOLLOW, "biden")
	fmt.Println(RelMgr("qing"))
	
	RelAction("qing", DO_BLOCK, "biden")
	fmt.Println(RelMgr("qing"))
}
