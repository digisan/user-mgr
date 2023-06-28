package registered

import (
	"fmt"
	"testing"

	. "github.com/digisan/user-mgr/user/tool"
)

func TestAdminFields(t *testing.T) {
	fmt.Println(ListField(Admin{}))

	for _, f := range ListField(Admin{}) {
		fmt.Printf("%s\n", f)
	}
}
