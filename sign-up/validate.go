package signup

import (
	"context"
	"fmt"
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
	if err := user.Validate(exclTags...); err != nil {
		return TransInvalidErr(user, err)
	}
	return nil
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
	code, err := verifyEmail(user)
	if err != nil {
		return err
	}
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
		return nil, fmt.Errorf("no email code exists")
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
