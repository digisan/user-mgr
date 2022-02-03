package signup

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/tool"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	"gopkg.in/go-playground/validator.v9"
)

const (
	MinLenLetter = 6
	timeout      = 30 * time.Second
)

var (
	mUserCodeTm = &sync.Map{}
)

func validateInput(user usr.User) error {

	v := validator.New()
	_ = v.RegisterValidation("active", func(fl validator.FieldLevel) bool {
		active := fl.Field().String()
		return usr.VerifyActive(active)
	})
	_ = v.RegisterValidation("uname", func(fl validator.FieldLevel) bool {
		uname := fl.Field().String()
		return !udb.UserDB.IsExisting(uname, false)
	})
	_ = v.RegisterValidation("pwd", func(fl validator.FieldLevel) bool {
		lenOK, number, upper, special := usr.VerifyPwd(fl.Field().String(), MinLenLetter)
		return lenOK && number && upper && special
	})
	_ = v.RegisterValidation("tel", func(fl validator.FieldLevel) bool {
		tel := fl.Field().String()
		return usr.VerifyTel(tel)
	})
	_ = v.RegisterValidation("addr", func(fl validator.FieldLevel) bool {
		addr := fl.Field().String()
		return usr.VerifyAddr(addr)
	})
	_ = v.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		role := fl.Field().String()
		return usr.VerifyAddr(role)
	})
	_ = v.RegisterValidation("level", func(fl validator.FieldLevel) bool {
		level := fl.Field().String()
		return usr.VerifyAddr(level)
	})
	_ = v.RegisterValidation("expire", func(fl validator.FieldLevel) bool {
		expire := fl.Field().String()
		return usr.VerifyAddr(expire)
	})
	_ = v.RegisterValidation("avatar", func(fl validator.FieldLevel) bool {
		avatar := fl.Field().String()
		return usr.VerifyAddr(avatar)
	})

	if err := v.Struct(user); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			lk.WarnOnErr("%v", e)
			es := fmt.Sprint(e)
			switch {
			case strings.Contains(es, "Active"):
				return fmt.Errorf("invalid active status, only set 'T' for true, or 'F' for false")
			case strings.Contains(es, "Password"):
				return fmt.Errorf("invalid password, at least %d letters, consist of UPPER-CASE, number and symbol", MinLenLetter)
			case strings.Contains(es, "UName"):
				return fmt.Errorf("invalid user name, [%s] is already existing", user.UName)
			case strings.Contains(es, "Tel"):
				return fmt.Errorf("invalid telephone number")
			}
		}
	}

	return nil
}

func verifyEmail(user usr.User) (string, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	code, err := tool.SendCode(ctx, user.Email, timeout)
	if err != nil {
		return "", fmt.Errorf("verification code sending error: %v", err)
	}

	return code, nil
}

// POST 1
func Validate(user usr.User, chkInput, chkEmail bool) error {

	if chkInput {
		if err := validateInput(user); err != nil {
			return err
		}
	}

	if chkEmail {
		var (
			code string
			err  error
		)
		if code, err = verifyEmail(user); err != nil {
			return err
		}
		mUserCodeTm.Store(user.UName, struct {
			Code string
			Tm   time.Time
		}{code, time.Now()})
	}

	return nil
}

// POST 2
func VerifyCode(user usr.User, incode string) error {

	val, ok := mUserCodeTm.LoadAndDelete(user.UName)
	if !ok {
		return fmt.Errorf("no email code exists")
	}
	ct := val.(struct {
		Code string
		Tm   time.Time
	})

	if time.Since(ct.Tm) > 30*time.Minute {
		return fmt.Errorf("email verification code is expired")
	}

	// fmt.Println("Input your code sent to you email")
	// incode := ""
	// fmt.Scanf("%s", &incode)

	if ct.Code != incode {
		return fmt.Errorf("email couldn't be verified")
	}

	return nil
}

// only support gmail now
func SetCodeEmail(email, password string) {
	lk.FailOnErr("%v", tool.SetGMail(email, password))
}
