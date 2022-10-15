package signup

import (
	"fmt"
	"sync"
	"time"

	"github.com/digisan/user-mgr/tool"
	usr "github.com/digisan/user-mgr/user"
)

var (
	timeoutVerify = 10 * time.Minute
	mUserCodeTm   = &sync.Map{}
)

func SetVerifyEmailTimeout(t time.Duration) {
	timeoutVerify = t
}

// POST 1
func ChkInput(user *usr.User, exclTags ...string) error {
	return user.Validate(exclTags...)
}

func verifyEmail(user *usr.User) (string, error) {
	return tool.SendCode(user.Email)
}

// POST 1
func ChkEmail(user *usr.User) error {

	var (
		code string
		err  error
	)

	// backdoor for debugging
	// {
	// 	if strings.HasPrefix(user.Password, "*") {
	// 		user.MemLevel = 3 // admin, 2 advanced, 1 subscribe, 0 registered
	// 		code = user.Password
	// 		goto STORE
	// 	}
	// }

	code, err = verifyEmail(user)
	if err != nil {
		return err
	}

	// STORE:
	mUserCodeTm.Store(user.UName, struct {
		Code string
		Tm   time.Time
		user *usr.User
	}{code, time.Now(), user})
	return nil
}

// POST 2
func VerifyCode(uname, incode string) (*usr.User, error) {

	// fmt.Println("Input your code sent to you email")
	// incode := ""
	// fmt.Scanf("%s", &incode)

	val, ok := mUserCodeTm.Load(uname)
	if !ok {
		return nil, fmt.Errorf("there is no email verification code for [%s]", uname)
	}

	ctu := val.(struct {
		Code string
		Tm   time.Time
		user *usr.User
	})

	if time.Since(ctu.Tm) > timeoutVerify {
		return nil, fmt.Errorf("email verification code is expired")
	}

	if ctu.Code != incode {
		return nil, fmt.Errorf("email couldn't be verified")
	}

	mUserCodeTm.Delete(uname)
	return ctu.user, nil
}
