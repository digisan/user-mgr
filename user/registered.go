package user

import (
	"errors"
	"fmt"

	bh "github.com/digisan/db-helper/badger"
	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/user/db"
	. "github.com/digisan/user-mgr/user/registered"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

///////////////////////////////////////////////////

func RemoveUser(uname string, lock bool) error {
	if lock {
		db.DbGrp.Lock()
		defer db.DbGrp.Unlock()
	}
	prefixes := [][]byte{
		[]byte("T" + SEP + uname + SEP),
		[]byte("F" + SEP + uname + SEP),
	}
	for _, prefix := range prefixes {
		n, err := bh.DeleteFirstObject[User](prefix)
		if err != nil {
			return err
		}
		if n == 1 {
			break
		}
	}
	return nil
}

func UpdateUser(u *User) error {
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	if err := RemoveUser(u.UName, false); err != nil {
		return err
	}
	return bh.UpsertOneObject(u)
}

func LoadUser(uname string, active bool) (*User, bool, error) {
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	prefix := []byte("T" + SEP + uname + SEP)
	if !active {
		prefix = []byte("F" + SEP + uname + SEP)
	}
	u, err := bh.GetFirstObject[User](prefix, nil)
	return u, err == nil && u != nil && u.Email != "", err
}

func LoadActiveUser(uname string) (*User, bool, error) {
	return LoadUser(uname, true)
}

func LoadAnyUser(uname string) (*User, bool, error) {
	uA, okA, errA := LoadUser(uname, true)
	uD, okD, errD := LoadUser(uname, false)
	var u *User
	if okA {
		u = uA
	} else if okD {
		u = uD
	}
	var err error
	if errA != nil {
		err = errA
	} else if errD != nil {
		err = errD
	}
	return u, err == nil && (okA || okD), err
}

func LoadUserByUniProp(propName, propVal string, active bool) (*User, bool, error) {
	var (
		err error
	)
	users, err := ListUser(func(u *User) bool {
		flag := u.IsActive()
		if !active {
			flag = !u.IsActive()
		}
		switch propName {
		case "uname", "Uname":
			return flag && u.UName == propVal
		case "email", "Email":
			return flag && u.Email == propVal
		case "phone", "Phone":
			return flag && u.Phone == propVal
		default:
			return false
		}
	})
	if len(users) > 0 {
		u := users[0]
		return u, err == nil && u != nil && u.Email != "", err
	}
	return nil, false, err
}

func LoadActiveUserByUniProp(propName, propVal string) (*User, bool, error) {
	return LoadUserByUniProp(propName, propVal, true)
}

func LoadAnyUserByUniProp(propName, propVal string) (*User, bool, error) {
	uA, okA, errA := LoadUserByUniProp(propName, propVal, true)
	uD, okD, errD := LoadUserByUniProp(propName, propVal, false)
	var u *User
	if okA {
		u = uA
	} else if okD {
		u = uD
	}
	var err error
	if errA != nil {
		err = errA
	} else if errD != nil {
		err = errD
	}
	return u, err == nil && (okA || okD), err
}

func ListUser(filter func(*User) bool) ([]*User, error) {
	db.DbGrp.Lock()
	defer db.DbGrp.Unlock()

	return bh.GetObjects([]byte(""), filter)
}

func UserExists(uname, email string, activeOnly bool) bool {
	if activeOnly {
		// check uname
		_, ok, err := LoadUser(uname, true)
		lk.WarnOnErr("%v", err)
		if ok {
			return ok
		}
		// check email
		_, ok, err = LoadActiveUserByUniProp("email", email)
		lk.WarnOnErr("%v", err)
		return ok

	} else {
		// check uname
		_, ok, err := LoadAnyUser(uname)
		lk.WarnOnErr("%v", err)
		if ok {
			return ok
		}
		// check email
		_, ok, err = LoadAnyUserByUniProp("email", email)
		lk.WarnOnErr("%v", err)
		return ok
	}
}

// only for unique value
func UsedByOther(uname_self, propName, propVal string) bool {
	u, ok, err := LoadAnyUserByUniProp(propName, propVal)
	if err == nil && ok && u != nil {
		return uname_self != u.UName
	}
	return false
}

func SetUserBoolField(uname, field string, flag bool) (u *User, ok bool, err error) {
	if u, ok, err = LoadAnyUser(uname); err == nil {
		if ok {
			switch field {
			case "Active", "active", "ACTIVE":
				u.Active = flag
			case "Official", "official", "OFFICIAL":
				u.Official = flag
			case "Certified", "certified", "CERTIFIED":
				u.Certified = flag
			default:
				lk.FailOnErr("%v", fmt.Errorf("[%s] is unsupported setting BoolField", field))
			}
			if err = UpdateUser(u); err != nil {
				return nil, false, err
			}
			u, ok, err = LoadAnyUser(uname)
			return u, err == nil && ok, err
		}
		return nil, false, fmt.Errorf("couldn't find [%s] for setting [%s]", uname, field)
	}
	return nil, false, err
}

func ActivateUser(uname string, flag bool) (*User, bool, error) {
	return SetUserBoolField(uname, "active", flag)
}

func OfficializeUser(uname string, flag bool) (*User, bool, error) {
	return SetUserBoolField(uname, "official", flag)
}

func CertifyUser(uname string, flag bool) (*User, bool, error) {
	return SetUserBoolField(uname, "certified", flag)
}

//////////////////////////////////////////////////////////////////

// to fetch field from "claims", map key must be json key.
// may not struct field name.
func TokenClaimsInHandler(c echo.Context) (*jwt.Token, jwt.MapClaims, error) {
	if c.Get("user") == nil {
		return nil, nil, errors.New("JWT token missing")
	}
	token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
	if !ok {
		return nil, nil, errors.New("JWT token invalid")
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

func ToFullUser(c echo.Context) (*User, error) {
	_, claims, err := TokenClaimsInHandler(c)
	if err != nil {
		return nil, err
	}
	userSlim := ClaimsToUser(claims)
	user, ok, err := LoadAnyUser(userSlim.UName)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("cannot find [%s] for its full fields", userSlim.UName)
	}
	return user, nil
}

func ToActiveFullUser(c echo.Context) (*User, error) {
	_, claims, err := TokenClaimsInHandler(c)
	if err != nil {
		return nil, err
	}
	userSlim := ClaimsToUser(claims)
	user, ok, err := LoadActiveUser(userSlim.UName)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("cannot find active [%s] for its full fields", userSlim.UName)
	}
	return user, nil
}
