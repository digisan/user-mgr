package user

import (
	"errors"
	"sync"
	"time"

	fd "github.com/digisan/gotk/filedir"
	lk "github.com/digisan/logkit"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type UserClaims struct {
	Core
	jwt.RegisteredClaims
}

var (
	key     = fd.SelfMD5()
	smToken = &sync.Map{}
)

func TokenKey() string {
	return key
}

// invoke in 'login'
func MakeClaims(user *User) *UserClaims {
	now := time.Now()
	return &UserClaims{
		user.Core,
		jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  []string{},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 72)),
			NotBefore: &jwt.NumericDate{},
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        "",
		},
	}
}

// invoke in 'login', store token in cache here
func GenerateToken(uc *UserClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	ts, err := token.SignedString([]byte(key))
	lk.FailOnErr("%v", err)
	smToken.Store(uc.UName, ts)
	return ts
}

// invoke in 'logout'
func (u *User) DeleteToken() {
	smToken.Delete(u.UName)
}

// validate token
func (u *User) ValidateToken(token string) bool {
	tkn, ok := smToken.Load(u.UName)
	return ok && tkn == token
}

// to fetch field from "claims", map key must be json key.
// may not struct field name.
func TokenClaimsInHandler(c echo.Context) (*jwt.Token, jwt.MapClaims, error) {
	token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
	if !ok {
		return nil, nil, errors.New("JWT token missing or invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return nil, nil, errors.New("failed to cast claims as jwt.MapClaims")
	}
	return token, claims, nil
}

func ClaimsToUser(claims jwt.MapClaims) *User {
	return &User{
		Core: Core{
			UName:    claims["uname"].(string),
			Email:    claims["email"].(string),
			Password: claims["password"].(string),
		},
		Profile: Profile{},
		Admin:   Admin{},
	}
}

func Invoker(c echo.Context) (*User, error) {
	_, claims, err := TokenClaimsInHandler(c)
	if err != nil {
		return nil, err
	}
	return ClaimsToUser(claims), nil
}
