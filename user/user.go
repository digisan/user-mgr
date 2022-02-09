package user

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/digisan/go-generics/str"
	"github.com/digisan/user-mgr/tool"
)

type User struct {
	Active     string `json:"active" validate:"active"`         // "T" "F"
	UName      string `json:"uname" validate:"required,uname"`  // unique, registered name
	Email      string `json:"email" validate:"required,email"`  // unique
	Name       string `json:"name" validate:"required,name"`    // real name
	Password   string `json:"password" validate:"required,pwd"` // <-- a custom validation rule, plaintext!
	Regtime    string `json:"regtime" validate:"regtime"`       // register time
	Phone      string `json:"phone" validate:"phone"`           // optional
	Addr       string `json:"addr" validate:"addr"`             // optional
	SysRole    string `json:"role" validate:"sysrole"`          // optional
	SysLevel   string `json:"level" validate:"syslevel"`        // optional
	Expire     string `json:"expire" validate:"expire"`         // optional
	NationalID string `json:"nationalid" validate:"nationalid"` // optional
	Gender     string `json:"gender" validate:"gender"`         // optional
	Position   string `json:"position" validate:"position"`     // optional
	Title      string `json:"title" validate:"title"`           // optional
	Employer   string `json:"employer" validate:"employer"`     // optional
	Avatar     string `json:"avatar" validate:"avatar"`         // optional
	key        []byte // at last, from 'Password'
}

func ListUserField() (fields []string) {
	typ := reflect.TypeOf(User{})
	// fmt.Println("Type:", typ.Name(), "Kind:", typ.Kind())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fields = append(fields, field.Name)
	}
	return
}

func ListUserValidator() (tags []string) {
	typ := reflect.TypeOf(User{})
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("validate")
		// fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
		tags = append(tags, strings.Split(tag, ",")...)
	}
	return str.FM(str.MkSet(tags...),
		func(i int, e string) bool {
			return len(e) > 0 && str.NotIn(e, "required", "email") // exclude internal validate tags
		},
		nil,
	)
}

func (u User) String() string {
	if u.UName != "" {
		return fmt.Sprintf(
			"Active: %s\n"+
				"UName: %s\n"+
				"Email: %s\n"+
				"Name: %s\n"+
				"Password: %s\n"+
				"Register Time: %s\n"+
				"Phone: %s\n"+
				"Address: %s\n"+
				"Role: %s\n"+
				"Level: %s\n"+
				"Expire: %s\n"+
				"Avatar: %s", // add more
			u.Active, u.UName, u.Email,
			u.Name, u.Password, u.Regtime,
			u.Phone, u.Addr, u.SysRole,
			u.SysLevel, u.Expire, u.Avatar,
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

// Active||UName||Email||Name||Regtime||Phone||...||pwdBuf <==> key
func (u *User) Marshal() (info []byte, key []byte) {
	// key : db value
	key = u.GenKey() // db value
	// info : db key
	info = []byte(fmt.Sprintf("%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||",
		u.Active, u.UName, u.Email,
		u.Name, u.Regtime, u.Phone,
		u.Addr, u.SysRole, u.SysLevel,
		u.Expire, u.Avatar)) // add more
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
		case 5: // Phone
			u.Phone = value
		case 6: // Address
			u.Addr = value
		case 7: // Role
			u.SysRole = value
		case 8: // Level
			u.SysLevel = value
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

func (u *User) IsActive() bool {
	return u.Active == "T"
}

func (u *User) StampRegTime() {
	u.Regtime = time.Now().UTC().Format(time.RFC3339)
}
