package registered

import (
	"fmt"
	"testing"

	. "github.com/digisan/user-mgr/user/tool"
)

func TestCoreFields(t *testing.T) {
	fmt.Println(ListField(Core{}))

	for _, f := range ListField(Core{}) {
		fmt.Printf("%s\n", f)
	}
}
