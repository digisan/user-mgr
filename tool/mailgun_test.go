package tool

import (
	"fmt"
	"testing"
)

func TestMail(t *testing.T) {
	ok, sent, failed, errs := SendMail("Fancy subject!", "Hello from Mailgun Go!", "cdutwhu@outlook.com", "4987346@qq.com")

	fmt.Println(ok)
	if ok {
		fmt.Println(sent)
		fmt.Println("---")
		fmt.Println(failed)
		fmt.Println("---")
		fmt.Println(errs)
	}	
}
