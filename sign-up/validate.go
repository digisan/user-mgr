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
	timeoutVerify = 30 * time.Minute
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
func ChkInput(user usr.User) error {
	if err := user.Validate(); err != nil {
		return TransInvalidErr(err)
	}
	return nil
}

func verifyEmail(user usr.User) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	code, err := tool.SendCode(ctx, user.Email, timeoutSend)
	if err != nil {
		return "", fmt.Errorf("verification code sending error: %v", err)
	}
	return code, nil
}

// POST 1
func ChkEmail(user usr.User) error {
	code, err := verifyEmail(user)
	if err != nil {
		return err
	}
	mUserCodeTm.Store(user.UName, struct {
		Code string
		Tm   time.Time
	}{code, time.Now()})
	return nil
}

// POST 2
func VerifyCode(user usr.User, incode string) error {

	// fmt.Println("Input your code sent to you email")
	// incode := ""
	// fmt.Scanf("%s", &incode)

	val, ok := mUserCodeTm.LoadAndDelete(user.UName)
	if !ok {
		return fmt.Errorf("no email code exists")
	}

	ct := val.(struct {
		Code string
		Tm   time.Time
	})

	if time.Since(ct.Tm) > timeoutVerify {
		return fmt.Errorf("email verification code is expired")
	}

	if ct.Code != incode {
		return fmt.Errorf("email couldn't be verified")
	}

	return nil
}
