package signup

import (
	"sync"
	"time"

	. "github.com/digisan/user-mgr/cst"
	u "github.com/digisan/user-mgr/user"
	ur "github.com/digisan/user-mgr/user/registered"
	. "github.com/digisan/user-mgr/util"
)

var (
	timeoutVerify = 10 * time.Minute
	mUserCodeTm   = &sync.Map{}
)

func SetVerifyEmailTimeout(t time.Duration) {
	timeoutVerify = t
}

// POST 1
func CheckInput(user *ur.User, exclTags ...string) error {
	return u.Validate(user, exclTags...)
}

func verifyEmail(user *ur.User) (string, error) {
	return SendCode(user.Email)
}

// POST 1
func CheckEmail(user *ur.User) error {

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
		user *ur.User
	}{code, time.Now(), user})
	return nil
}

// POST 2
func VerifyCode(uname, incode string) (*ur.User, error) {

	// fmt.Println("Input your code sent to you email")
	// incode := ""
	// fmt.Scanf("%s", &incode)

	val, ok := mUserCodeTm.Load(uname)
	if !ok {
		return nil, Err(ERR_VCODE_MISSING).Wrap(uname)
	}

	ctu := val.(struct {
		Code string
		Tm   time.Time
		user *ur.User
	})

	if time.Since(ctu.Tm) > timeoutVerify {
		return nil, Err(ERR_VCODE_EXP)
	}

	if ctu.Code != incode {
		return nil, Err(ERR_VCODE_VERIFY_FAIL).Wrap("email")
	}

	mUserCodeTm.Delete(uname)
	return ctu.user, nil
}
