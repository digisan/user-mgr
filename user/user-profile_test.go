package user

import (
	"fmt"
	"testing"
)

func TestProfileFields(t *testing.T) {
	fmt.Println(ListField(Profile{}))

	for _, f := range ListField(Profile{}) {
		fmt.Printf("%s\n", f)
	}
}
