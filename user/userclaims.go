package user

import (
	"sync"
	"time"

	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/tool"
	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	User
	jwt.StandardClaims
}

var (
	key        = tool.SelfMD5()
	mUserToken = &sync.Map{}
)

func TokenKey() string {
	return key
}

// invoke in 'login'
func MakeUserClaims(user User) *UserClaims {
	return &UserClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
}

// invoke in 'login'
func (uc *UserClaims) GenToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	ts, err := token.SignedString([]byte(key))
	lk.FailOnErr("%v", err)
	mUserToken.Store(uc.UName, ts)
	return ts
}

// invoke in 'logout'
func (uc *UserClaims) DeleteToken() {
	mUserToken.Delete(uc.UName)
}

// validate token
func (uc *UserClaims) ValidateToken(token string) bool {
	tkn, ok := mUserToken.Load(uc.UName)
	return ok && tkn == token
}
