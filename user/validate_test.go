package user

import (
	"fmt"
	"testing"

	ur "github.com/digisan/user-mgr/user/registered"
	. "github.com/digisan/user-mgr/util"
)

func TestListValidator(t *testing.T) {
	vTags := ListValidator(ur.User{}.Core, ur.User{}.Profile, ur.User{}.Admin)
	fmt.Println(vTags)

	u := ur.User{
		Core:    ur.Core{},
		Profile: ur.Profile{},
		Admin:   ur.Admin{},
	}
	vTags = ListValidator(u.Core, u.Profile, u.Admin)
	fmt.Println(vTags)
}
