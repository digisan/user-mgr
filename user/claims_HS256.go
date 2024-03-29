package user

// import (
// 	"context"
// 	"errors"
// 	"sync"
// 	"time"

// 	fd "github.com/digisan/gotk/file-dir"
// 	lk "github.com/digisan/logkit"
// 	ur "github.com/digisan/user-mgr/user/registered"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/labstack/echo/v4"
// )

// type UserClaims struct {
// 	ur.Core
// 	jwt.RegisteredClaims
// }

// type TokenInfo struct {
// 	value string
// 	start time.Time
// }

// var (
// 	key         = fd.SelfMD5()
// 	smToken     = &sync.Map{}    // uname: *TokenInfo
// 	validPeriod = time.Hour * 24 // default token valid period
// )

// func TokenKey() []byte {
// 	return []byte(key)
// }

// func MonitorTokenExpired(ctx context.Context, fnOnGotTokenExp func(uname string) error) {
// 	const interval = 15 * time.Second
// 	go func(ctx context.Context) {
// 		ticker := time.NewTicker(interval)
// 		for {
// 			select {
// 			case <-ticker.C:
// 				expUsers := []string{}
// 				smToken.Range(func(key, value any) bool {
// 					uname := key.(string)
// 					tkInfo := value.(*TokenInfo)
// 					if time.Since(tkInfo.start) > validPeriod {
// 						expUsers = append(expUsers, uname)
// 						if fnOnGotTokenExp != nil {
// 							lk.WarnOnErr("%v", fnOnGotTokenExp(uname))
// 						}
// 					}
// 					return true
// 				})
// 				for _, user := range expUsers {
// 					smToken.Delete(user)
// 				}
// 			case <-ctx.Done():
// 				return
// 			}
// 		}
// 	}(ctx)
// }

// // must invoke this before 'MakeClaims'
// func SetTokenValidPeriod(period time.Duration) {
// 	validPeriod = period
// }

// // invoke in 'login'
// func MakeClaims(user *ur.User) *UserClaims {
// 	now := time.Now()
// 	return &UserClaims{
// 		user.Core,
// 		jwt.RegisteredClaims{
// 			Issuer:    "",
// 			Subject:   "",
// 			Audience:  []string{},
// 			ExpiresAt: jwt.NewNumericDate(now.Add(validPeriod)),
// 			NotBefore: &jwt.NumericDate{},
// 			IssuedAt:  jwt.NewNumericDate(now),
// 			ID:        "",
// 		},
// 	}
// }

// // invoke in 'login', store token in cache here
// func GenerateToken(uc *UserClaims) string {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
// 	ts, err := token.SignedString(TokenKey())
// 	lk.FailOnErr("%v", err)
// 	smToken.Store(uc.UName, &TokenInfo{
// 		value: ts,
// 		start: time.Now(),
// 	})
// 	return ts
// }

// // invoke in 'logout'
// func DeleteToken(u *ur.User) {
// 	smToken.Delete(u.UName)
// }

// // validate token
// func ValidateToken(u *ur.User, token string) bool {
// 	tkInfo, ok := smToken.Load(u.UName)
// 	return ok && tkInfo.(*TokenInfo).value == token
// }

// // to fetch field from "claims", map key must be json key.
// // may not struct field name.
// func TokenClaimsInHandler(c echo.Context) (*jwt.Token, jwt.MapClaims, error) {
// 	token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
// 	if !ok {
// 		return nil, nil, errors.New("JWT token missing or invalid")
// 	}
// 	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
// 	if !ok {
// 		return nil, nil, errors.New("failed to cast claims as jwt.MapClaims")
// 	}
// 	return token, claims, nil
// }

// func ClaimsToUser(claims jwt.MapClaims) *ur.User {
// 	return &ur.User{
// 		Core: ur.Core{
// 			UName:    claims["uname"].(string),
// 			Email:    claims["email"].(string),
// 			Password: claims["password"].(string),
// 		},
// 		Profile: ur.Profile{},
// 		Admin:   ur.Admin{},
// 	}
// }

// func Invoker(c echo.Context) (*ur.User, error) {
// 	_, claims, err := TokenClaimsInHandler(c)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ClaimsToUser(claims), nil
// }
