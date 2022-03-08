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
	// RelAction(DO_FOLLOW, "qing", "trump")
	// RelAction(DO_UNFOLLOW, "qing", "musk")
	// RelAction(DO_FOLLOW, "qing", "musk")

	relQing := RelMgr(FOLLOWING, "qing")
	fmt.Println(relQing)

	fmt.Print("-----------------------")

	relMusk := RelMgr(FOLLOWER, "trump")
	fmt.Println(relMusk)

}
