package user

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/digisan/go-generics/str"
	"github.com/digisan/user-mgr/tool"
)

// if modified, change 1. MOK_***, 2. mFldAddr, 3. 'auto-tags.go', 4. 'validator.go' in sign-up.
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
	Tags       string `json:"tags" validate:"tags"`             // optional // linked by '^'
	Avatar     []byte `json:"avatar" validate:"avatar"`         // optional
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
		t, v := reflect.TypeOf(u), reflect.ValueOf(u)
		for i := 0; i < t.NumField(); i++ {
			fld, val := t.Field(i), v.Field(i)
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

const (
	SEP     = "||"
	SEP_TAG = "^"
)

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
	MOK_Tags
	MOK_Avatar
	MOK_PwdBuf
	MOK_END
)

func (u *User) KeyFieldAddr(mok int) interface{} {
	mFldAddr := map[int]interface{}{
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
		MOK_Tags:       &u.Tags,
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

func (u *User) ValFieldAddr(mov int) interface{} {
	mFldAddr := map[int]interface{}{
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
		fnFldAddr func(int) interface{}
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
			switch v := param.fnFldAddr(i).(type) {
			case *string:
				sb.Write([]byte(*v))
			case *[]byte:
				sb.Write(*v)
			default:
				panic("Marshal Error Type")
			}
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
		fnFldAddr func(int) interface{}
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
			switch v := param.fnFldAddr(i).(type) {
			case *string:
				*v = string(seg)
			case *[]byte:
				*v = seg
			default:
				panic("Unmarshal Error Type")
			}
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

func (u *User) GetTags() []string {
	return strings.Split(u.Tags, SEP_TAG)
}

func (u *User) AddTags(tags ...string) {
	tagsExs := strings.Split(u.Tags, SEP_TAG)
	tags = append(tags, tagsExs...)
	tags = str.MkSet(tags...)
	u.Tags = strings.TrimSuffix(strings.Join(tags, SEP_TAG), SEP_TAG)
}

func (u *User) RmTags(tags ...string) {
	tagsExs := strings.Split(u.Tags, SEP_TAG)
	tags = str.Minus(tagsExs, tags)
	u.Tags = strings.TrimSuffix(strings.Join(tags, SEP_TAG), SEP_TAG)
}

func (u *User) SetAvatar(r io.Reader) {
	u.Avatar = tool.StreamToByte(r)
}

func (u *User) AvatarStdBase64() string {
	return base64.StdEncoding.EncodeToString(u.Avatar)
}

func (u *User) AvatarUrlBase64() string {
	return base64.URLEncoding.EncodeToString(u.Avatar)
}
