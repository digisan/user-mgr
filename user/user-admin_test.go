package user

import (
	"fmt"
	"testing"
)

func TestAdminFields(t *testing.T) {
	fmt.Println(ListField(Admin{}))

	for _, f := range ListField(Admin{}) {
		fmt.Printf("%s\n", f)
	}
}
