package user

import (
	"fmt"
	"testing"
)

func TestCoreFields(t *testing.T) {
	fmt.Println(ListField(Core{}))

	for _, f := range ListField(Core{}) {
		fmt.Printf("%s\n", f)
	}
}
