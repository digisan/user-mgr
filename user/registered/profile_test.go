package registered

import (
	"fmt"
	"testing"

	. "github.com/digisan/user-mgr/util"
)

func TestProfileFields(t *testing.T) {
	fmt.Println(ListField(Profile{}))

	for _, f := range ListField(Profile{}) {
		fmt.Printf("%s\n", f)
	}
}
