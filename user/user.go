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
	SysRole    string `json:"role" validate:"sysRole"`          // optional
	MemLevel   string `json:"level" validate:"memLevel"`        // optional
	MemExpire  string `json:"expire" validate:"memExpire"`      // optional
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
		sb := strings.Builder{}
		typ := reflect.TypeOf(u)
		val := reflect.ValueOf(u)
		for i := 0; i < typ.NumField(); i++ {
			fld := typ.Field(i)
			val := val.Field(i)
			sb.WriteString(fmt.Sprintf("%-10s %v\n", fld.Name+":", val.String()))
		}
		return sb.String()
	}
	return "[Empty User]\n"
}

func (u *User) GenKey() []byte {
	if u.key == nil {
		u.key = []byte(fmt.Sprintf("%d", time.Now().UnixNano())[3:19])
	}
	return u.key
}

// Active||UName||Email||Name||...||pwdBuf <==> key
func (u *User) Marshal() (info []byte, key []byte) {
	// key : db value
	key = u.GenKey() // db value
	// info : db key
	info = []byte(fmt.Sprintf("%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||%s||",
		u.Active, u.UName, u.Email,
		u.Name, u.Regtime, u.Phone,
		u.Addr, u.SysRole, u.MemLevel,
		u.MemExpire, u.NationalID, u.Gender,
		u.Position, u.Title, u.Employer,
		u.Avatar)) // add more
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
			u.MemLevel = value
		case 9: // Expire
			u.MemExpire = value
		case 10: // National ID
			u.NationalID = value
		case 11: // Gender
			u.Gender = value
		case 12: // Position
			u.Position = value
		case 13: // Title
			u.Title = value
		case 14: // Employer
			u.Employer = value
		case 15: // Avatar
			u.Avatar = value
		// add more
		case 16: // pwdBuf (16 must change to be the last one if added more)
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
