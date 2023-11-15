package util

import (
	"fmt"
	"time"

	gm "github.com/digisan/go-mail"
	"github.com/digisan/gotk/crypto"
	"github.com/digisan/gotk/strs"
	. "github.com/digisan/user-mgr/cst"
)

func genCode(email string) string {
	key := []byte(fmt.Sprintf("%d", time.Now().UnixNano())[3:19])
	return strs.MaxLen(fmt.Sprintf("%06x", crypto.Encrypt(email, key)), 6)
}

// obsoleted !!!
// func SendCode(ctx context.Context, recipient string, timeout time.Duration) (string, error) {

// 	header := "Sign-Up Verification Code"

// 	code := genCode(recipient)
// 	// fmt.Println(code)
// 	body := fmt.Sprintf("verification code: %s\n", code)

// 	chOK := make(chan error)
// 	go gmail(chOK, from, pwd, header, body, recipient)

// 	select {
// 	case <-ctx.Done():
// 		// fmt.Printf("%v\n", "out cancelled")
// 		return "", fmt.Errorf("out cancelled")

// 	case <-time.After(timeout):
// 		// fmt.Printf("%v\n", "time out")
// 		return "", fmt.Errorf("time out")

// 	case err := <-chOK:
// 		// fmt.Printf("%v\n", err)
// 		return code, err
// 	}
// }

// if do unit test, comment out 'gm' invoked !!!
func SendCode(recipient string) (string, error) {

	subject := "Sign-Up Verification Code"

	code := genCode(recipient)
	// fmt.Println(code)
	body := fmt.Sprintf("verification code: %s\n", code)

	if ok, _, _, errs := gm.SendMail(subject, body, recipient); ok {
		return code, nil
	} else {
		if len(errs) > 0 {
			return "", errs[0]
		}
		return "", Err(ERR_SEND_EMAIL).Wrap(recipient)
	}
}
