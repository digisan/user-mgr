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
	key        string
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
			sb.WriteString(fmt.Sprintf("%-12s %v\n", fld.Name+":", val.String()))
		}
		return sb.String()
	}
	return "[Empty User]\n"
}

func (u *User) GenKey() string {
	if u.key == "" {
		u.key = fmt.Sprintf("%d", time.Now().UnixNano())[3:19]
	}
	return u.key
}

const SEP = "||"

// db key order
const (
	MOK_Active int = iota
	MOK_UName
	MOK_Email
	MOK_Name
	MOK_Regtime
	MOK_Phone
	MOK_Addr
	MOK_SysRole
	MOK_MemLevel
	MOK_MemExpire
	MOK_NationalID
	MOK_Gender
	MOK_Position
	MOK_Title
	MOK_Employer
	MOK_Avatar
	MOK_PwdBuf
	MOK_END
)

func (u *User) KeyFieldAddr(mok int) *string {
	mFldAddr := map[int]*string{
		MOK_Active:     &u.Active,
		MOK_UName:      &u.UName,
		MOK_Email:      &u.Email,
		MOK_Name:       &u.Name,
		MOK_Regtime:    &u.Regtime,
		MOK_Phone:      &u.Phone,
		MOK_Addr:       &u.Addr,
		MOK_SysRole:    &u.SysRole,
		MOK_MemLevel:   &u.MemLevel,
		MOK_MemExpire:  &u.MemExpire,
		MOK_NationalID: &u.NationalID,
		MOK_Gender:     &u.Gender,
		MOK_Position:   &u.Position,
		MOK_Title:      &u.Title,
		MOK_Employer:   &u.Employer,
		MOK_Avatar:     &u.Avatar,
		MOK_PwdBuf:     &u.Password,
	}
	return mFldAddr[mok]
}

// db value order
const (
	MOV_Key int = iota
	MOV_END
)

func (u *User) ValFieldAddr(mov int) *string {
	mFldAddr := map[int]*string{
		MOV_Key: &u.key,
	}
	return mFldAddr[mov]
}

////////////////////////////////////////////////////

func (u *User) Marshal() (forKey, forValue []byte) {

	key := u.GenKey()
	forValue = []byte(fmt.Sprint(len(u.UName)) + key) // *** fake key forValue in db ***

	params := []struct {
		end       int
		fnFldAddr func(int) *string
		out       *[]byte
	}{
		{
			end:       MOK_END,
			fnFldAddr: u.KeyFieldAddr,
			out:       &forKey,
		},
	}
	for _, param := range params {
		sb := &strings.Builder{}
		for i := 0; i < param.end; i++ {
			if i > 0 {
				sb.WriteString(SEP)
			}
			if i == MOK_PwdBuf {
				sb.Write(tool.Encrypt(u.Password, []byte(key))) // from u.Password
				continue
			}
			sb.WriteString(*param.fnFldAddr(i))
		}
		*param.out = []byte(sb.String())
	}
	return
}

func (u *User) Unmarshal(dbKey, dbVal []byte) {
	if dbVal != nil {
		u.key = string(dbVal) // *** fake key ***
	}
	params := []struct {
		in        []byte
		fnFldAddr func(int) *string
	}{
		{
			in:        dbKey,
			fnFldAddr: u.KeyFieldAddr,
		},
	}
	for _, param := range params {
		for i, seg := range bytes.Split(param.in, []byte(SEP)) {
			if i == MOK_PwdBuf {
				if u.key != "" {
					offset := len(fmt.Sprint(len(u.UName)))
					u.key = u.key[offset:]
					u.Password = tool.Decrypt(seg, []byte(u.key))
					continue
				}
			}
			*param.fnFldAddr(i) = string(seg)
		}
	}
}

///////////////////////////////////////////////////

func (u *User) IsActive() bool {
	return u.Active == "T"
}

func (u *User) StampRegTime() {
	u.Regtime = time.Now().UTC().Format(time.RFC3339)
}
