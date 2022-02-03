package user

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/digisan/user-mgr/tool"
)

type User struct {
	Active   string `json:"active" validate:"required,active"`   // "T" "F"
	UName    string `json:"uname" validate:"required,uname"`     // unique
	Email    string `json:"email" validate:"required,email"`     // unique
	Name     string `json:"name" validate:"required"`            // real name
	Password string `json:"password" validate:"required,pwd"`    // <-- a custom validation rule, plaintext!
	Regtime  string `json:"regtime" validate:"required,regtime"` // register time
	Tel      string `json:"tel" validate:"tel"`                  // optional
	Addr     string `json:"addr" validate:"addr"`                // optional
	Role     string `json:"role" validate:"role"`                // optional
	Level    string `json:"level" validate:"level"`              // optional
	Expire   string `json:"expire" validate:"expire"`            // optional
	Avatar   string `json:"avatar" validate:"avatar"`            // optional
	key      []byte // at last, from 'Password'
}

func (u User) String() string {
	if u.UName != "" {
		return fmt.Sprintf("Active: %s\n"+
			"UName: %s\n"+
			"Email: %s\n"+
			"Name: %s\n"+
			"Password: %s\n"+
			"Register Time: %s\n"+
			"Telephone: %s\n"+
			"Address: %s\n"+
			"Role: %s\n"+
			"Level: %s\n"+
			"Expire: %s\n"+
			"Avatar: %s", // add more
			u.Active, u.UName, u.Email, u.Name, u.Password, u.Regtime,
			u.Tel, u.Addr, u.Role, u.Level, u.Expire, u.Avatar,
		)
	}
	return "[Empty User]\n"
}

func (u *User) GenKey() []byte {
	if u.key == nil {
		u.key = []byte(fmt.Sprintf("%d", time.Now().UnixNano())[3:19])
	}
	return u.key
}

// Active||UName||Email||Name||Tel||...||pwdBuf <==> key
func (u *User) Marshal() (info []byte, key []byte) {
	// key : db value
	key = u.GenKey() // db value
	// info : db key
	info = []byte(fmt.Sprintf("%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||",
		u.Active, u.UName, u.Email, u.Name, u.Regtime, u.Tel,
		u.Addr, u.Role, u.Level, u.Expire, u.Avatar)) // add more
	pwdBuf := tool.Encrypt(u.Password, key)
	info = append(info, pwdBuf...) // from u.Password
	return
}

func (u *User) Unmarshal(info []byte, key []byte) {
	if key != nil {
		u.key = key
	}
	for i, seg := range bytes.Split(info, []byte("||")) {
		value := string(seg)
		switch i {
		case 0: // Active
			u.Active = value
		case 1: // UName
			u.UName = value
		case 2: // Email
			u.Email = value
		case 3: // Name
			u.Name = value
		case 4: // Register Time
			u.Regtime = value
		case 5: // Tel
			u.Tel = value
		case 6: // Address
			u.Addr = value
		case 7: // Role
			u.Role = value
		case 8: // Level
			u.Level = value
		case 9: // Expire
			u.Expire = value
		case 10: // Avatar
			u.Avatar = value
		// add more
		case 11: // pwdBuf (11 must change as last if added more)
			if key != nil {
				u.Password = tool.Decrypt(seg, key)
			}
		default:
			panic("sep || error")
		}
	}
}

func (u *User) Activate(val bool) {
	u.Active = strings.ToUpper(fmt.Sprint(val))[:1]
}

func (u *User) IsActive() bool {
	return u.Active == "T"
}

func (u *User) SetTel(tel string) {
	u.Tel = tel
}

func (u *User) SetAddr(addr string) {
	u.Addr = addr
}

func (u *User) SetRole(role string) {
	u.Role = role
}

func (u *User) SetLevel(level string) {
	u.Level = level
}

func (u *User) SetExpire(expire string) {
	u.Expire = expire
}

func (u *User) SetAvatar(avatar string) {
	u.Avatar = avatar
}
