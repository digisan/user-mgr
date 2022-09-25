package signup

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/tool"
	usr "github.com/digisan/user-mgr/user"
)

var (
	timeoutSend   = 45 * time.Second
	timeoutVerify = 10 * time.Minute
	mUserCodeTm   = &sync.Map{}
)

func SetVerifyEmailTimeout(t time.Duration) {
	timeoutVerify = t
}

// only support gmail now
func SetCodeEmail(email, password string) {
	lk.FailOnErr("%v", tool.SetGMail(email, password))
}

// POST 1
func ChkInput(user *usr.User, exclTags ...string) error {
	return user.Validate(exclTags...)
}

func verifyEmail(user *usr.User) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	code, err := tool.SendCode(ctx, user.Email, timeoutSend)
	if err != nil {
		return "", fmt.Errorf("verification code sending error: %v", err)
	}
	return code, nil
}

// POST 1
func ChkEmail(user *usr.User) error {

	var (
		code string
		err  error
	)

	// backdoor for debugging
	{
		if strings.HasPrefix(user.Password, "*") {
			user.MemLevel = 3 // admin, 2 advanced, 1 subscribe, 0 registered
			code = user.Password
			goto STORE
		}
	}

	code, err = verifyEmail(user)
	if err != nil {
		return err
	}

STORE:
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
