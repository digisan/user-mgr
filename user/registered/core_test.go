package registered

import (
	"fmt"
	"testing"

	. "github.com/digisan/user-mgr/util"
)

func TestCoreFields(t *testing.T) {
	fmt.Println(ListField(Core{}))

	for _, f := range ListField(Core{}) {
		fmt.Printf("%s\n", f)
	}
}
