package registered

import (
	"fmt"
	"testing"

	. "github.com/digisan/user-mgr/util"
)

func TestAdminFields(t *testing.T) {
	fmt.Println(ListField(Admin{}))

	for _, f := range ListField(Admin{}) {
		fmt.Printf("%s\n", f)
	}
}
