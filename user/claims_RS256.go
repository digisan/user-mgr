package user

import (
	"context"
	"fmt"
	"sync"
	"time"

	lk "github.com/digisan/logkit"
	ur "github.com/digisan/user-mgr/user/registered"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	ur.Core
	jwt.RegisteredClaims
}

type TokenInfo struct {
	value string
	start time.Time
}

var (
	smToken     = &sync.Map{}    // uname: *TokenInfo
	validPeriod = time.Hour * 24 // default token valid period
)

// store a copy of token here for further validation
func (uc *UserClaims) GenerateToken(prvKey []byte) (string, error) {

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	// now := time.Now().UTC()
	// claims := make(jwt.MapClaims)
	// claims["dat"] = content             // Our custom data.
	// claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	// claims["iat"] = now.Unix()          // The time at which the token was issued.
	// claims["nbf"] = now.Unix()          // The time before which the token must be disregarded.

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, uc).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	smToken.Store(uc.UName, &TokenInfo{
		value: token,
		start: time.Now(),
	})

	return token, nil
}

func MonitorTokenExpired(ctx context.Context, cExpired chan<- string, fnOnGotTokenExp func(uname string) error) {
	const interval = 15 * time.Second
	go func(ctx context.Context) {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				expUsers := []string{}
				smToken.Range(func(key, value any) bool {
					uname := key.(string)
					tkInfo := value.(*TokenInfo)
					if time.Since(tkInfo.start) > validPeriod {
						expUsers = append(expUsers, uname)
						if fnOnGotTokenExp != nil {
							lk.WarnOnErr("%v", fnOnGotTokenExp(uname))
						}
						cExpired <- uname
					}
					return true
				})
				for _, user := range expUsers {
					smToken.Delete(user)
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}

// must invoke this before 'MakeClaims' !!!
func SetTokenValidPeriod(period time.Duration) {
	validPeriod = period
}

// invoke in 'login'
func MakeUserClaims(user *ur.User) *UserClaims {
	now := time.Now()
	return &UserClaims{
		user.Core,
		jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  []string{},
			ExpiresAt: jwt.NewNumericDate(now.Add(validPeriod)),
			NotBefore: &jwt.NumericDate{},
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        "",
		},
	}
}

// invoke in 'logout'
func DeleteToken(uname string) {
	smToken.Delete(uname)
}

// validate token
func ValidateToken(user *ur.User, ts string, pubKey []byte) (bool, error) {

	key, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return false, fmt.Errorf("validate: parse key: %w", err)
	}

	token, err := jwt.Parse(ts, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return false, fmt.Errorf("validate: %w", err)
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, fmt.Errorf("validate: token to MapClaims")
	}

	tkInfo, ok := smToken.Load(user.UName)
	if !ok || tkInfo.(*TokenInfo).value != ts {
		return false, fmt.Errorf("validate: token doesn't exist in record")
	}

	return true, nil
}
