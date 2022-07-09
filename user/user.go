package user

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
	bh "github.com/digisan/db-helper/badger-helper"
	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/crypto"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

// if modified, change 1. KO_***, 2. mFldAddr, 3. 'auto-tags.go', 4. 'validator.go' in sign-up.
type User struct {
	Core
	Profile
	Admin
}

func (u User) String() string {
	strCore := u.Core.String()
	if strings.Contains(strCore, "Empty") {
		return "[Empty User]\n"
	}
	return fmt.Sprintf("%s%s%s\n", strCore, u.Profile.String(), u.Admin.String())
}

// db key order
const (
	// Admin
	KO_Active int = iota
	// Core
	KO_UName
	KO_Email
	KO_Password
	KO_Key
	// Profile
	KO_Name
	KO_Phone
	KO_Country
	KO_City
	KO_Addr
	KO_PersonalIDType
	KO_PersonalID
	KO_Gender
	KO_DOB
	KO_Position
	KO_Title
	KO_Employer
	KO_Bio
	KO_AvatarType
	// Admin
	KO_Regtime
	KO_SysRole
	KO_MemLevel
	KO_MemExpire
	KO_Official
	KO_Certified
	KO_Tags
	//
	KO_END
)

// db value order
const (
	// Profile
	VO_Avatar int = iota
	//
	VO_END
)

func (u *User) KeyFieldAddr(ko int) any {
	mFldAddr := map[int]any{
		// Core
		KO_UName:    &u.UName,
		KO_Email:    &u.Email,
		KO_Password: &u.Password,
		KO_Key:      &u.key,
		// Profile
		KO_Name:           &u.Name,
		KO_Phone:          &u.Phone,
		KO_Country:        &u.Country,
		KO_City:           &u.City,
		KO_Addr:           &u.Addr,
		KO_PersonalIDType: &u.PersonalIDType,
		KO_PersonalID:     &u.PersonalID,
		KO_Gender:         &u.Gender,
		KO_DOB:            &u.DOB,
		KO_Position:       &u.Position,
		KO_Title:          &u.Title,
		KO_Employer:       &u.Employer,
		KO_Bio:            &u.Bio,
		KO_AvatarType:     &u.AvatarType,
		// Admin
		KO_Active:    &u.Active,
		KO_Regtime:   &u.Regtime,
		KO_SysRole:   &u.SysRole,
		KO_MemLevel:  &u.MemLevel,
		KO_MemExpire: &u.MemExpire,
		KO_Official:  &u.Official,
		KO_Certified: &u.Certified,
		KO_Tags:      &u.Tags,
	}
	return mFldAddr[ko]
}

func (u *User) ValFieldAddr(vo int) any {
	mFldAddr := map[int]any{
		// Profile
		VO_Avatar: &u.Avatar,
	}
	return mFldAddr[vo]
}

var secret = []int{
	KO_Email,
	KO_Password,
	KO_Name,
	KO_Phone,
	KO_Country,
	KO_City,
	KO_Addr,
	KO_PersonalIDType,
	KO_PersonalID,
	KO_Gender,
	KO_DOB,
	KO_Position,
	KO_Title,
	KO_Employer,
}

////////////////////////////////////////////////////

func (u *User) BadgerDB() *badger.DB {
	return DbGrp.Reg
}

func (u *User) Key() []byte {
	var (
		key = u.GenKey()
		sb  = &strings.Builder{}
	)
	for i := 0; i < KO_END; i++ {
		if i > 0 {
			sb.WriteString(SEP)
		}
		if In(i, secret...) {
			sb.Write(crypto.Encrypt((*u.KeyFieldAddr(i).(*string)), key[:]))
			continue
		}
		switch v := u.KeyFieldAddr(i).(type) {
		case *string:
			sb.WriteString(*v)
		case *[]byte:
			sb.Write(*v)
		case *[16]byte:
			sb.Write((*v)[:])
		case *bool:
			sb.WriteString(strings.ToUpper(fmt.Sprintf("%v", *v)[0:1]))
		case *time.Time:
			encoding, err := (*v).MarshalBinary()
			// lk.Debug(" --------------- encoding len: %d", len(encoding))
			lk.FailOnErr("%v", err)
			sb.Write(encoding)
		case *uint8:
			sb.Write([]byte{*v})
		default:
			panic("need more type for marshaling key")
		}
	}
	return []byte(sb.String())
}

func (u *User) Value() []byte {
	var (
		key = u.GenKey()
		sb  = &strings.Builder{}
	)
	for i := 0; i < VO_END; i++ {
		if i > 0 {
			sb.WriteString(SEP)
		}
		if In(i, secret...) {
			sb.Write(crypto.Encrypt((*u.ValFieldAddr(i).(*string)), key[:]))
			continue
		}
		switch v := u.ValFieldAddr(i).(type) {
		case *string:
			sb.WriteString(*v)
		case *[]byte:
			sb.Write(*v)
		default:
			panic("need more type for marshaling value")
		}
	}
	return []byte(sb.String())
}

func (u *User) Marshal(at any) (forKey, forValue []byte) {
	return u.Key(), u.Value()
}

func (u *User) Unmarshal(dbKey, dbVal []byte) (any, error) {
	params := []struct {
		in        []byte
		fnFldAddr func(int) any
	}{
		{dbKey, u.KeyFieldAddr},
		{dbVal, u.ValFieldAddr},
	}
	for ip, param := range params {
		if len(param.in) > 0 {

			var segs [][]byte

			if ip == 0 {
				segs = bytes.Split(param.in, []byte(SEP))
				u.key = *(*[16]byte)(segs[KO_Key])
			} else if ip == 1 {
				segs = [][]byte{param.in} // dbVal is one whole block
			}

			for i, seg := range segs {
				if (ip == 0 && i == KO_END) || (ip == 1 && i == VO_END) {
					break
				}
				if ip == 0 && In(i, secret...) {
					if u.key != [16]byte{} {
						*param.fnFldAddr(i).(*string) = crypto.Decrypt(seg, u.key[:])
						continue
					}
				}
				switch v := param.fnFldAddr(i).(type) {
				case *string:
					*v = string(seg)
				case *[]byte:
					*v = seg
				case *[16]byte:
					*v = *(*[16]byte)(seg)
				case *bool:
					if seg[0] == 'T' {
						*v = true
					} else {
						*v = false
					}
				case *time.Time:
					// lk.Debug(" --------------- seg len: %d", len(seg))
					t := &time.Time{}
					lk.FailOnErr("%v @ %v", t.UnmarshalBinary(seg), seg)
					*v = *t
				case *uint8:
					*v = seg[0]
				default:
					panic("Unmarshal Error Type")
				}
			}
		}
	}
	return u, nil
}

///////////////////////////////////////////////////

func (u *User) IsActive() bool {
	return u.Active
}

func (u *User) StampRegTime() {
	u.Regtime = time.Now()
}

func (u *User) SinceJoined() time.Duration {
	t := &time.Time{}
	t.UnmarshalText([]byte(u.Regtime.Format(time.RFC3339)))
	return time.Since(*t)
}

func (u *User) GetTags() []string {
	return strings.Split(u.Tags, SEP_TAG)
}

func (u *User) AddTags(tags ...string) {
	tagsExs := strings.Split(u.Tags, SEP_TAG)
	tags = append(tags, tagsExs...)
	tags = Settify(tags...)
	u.Tags = strings.TrimSuffix(strings.Join(tags, SEP_TAG), SEP_TAG)
}

func (u *User) RmTags(tags ...string) {
	tagsExs := strings.Split(u.Tags, SEP_TAG)
	tags = Minus(tagsExs, tags)
	u.Tags = strings.TrimSuffix(strings.Join(tags, SEP_TAG), SEP_TAG)
}

func (u *User) SetAvatar(avatarType string, r io.Reader) {
	u.AvatarType = avatarType
	u.Avatar = gio.StreamToBytes(r)
}

// 'avatarType' --- e.g. image/png, 'fh' --- FormFile('param')
// for example '<img src="data:image/png;base64,******/>'
func (u *User) SetAvatarByFormFile(avatarType string, fh *multipart.FileHeader) error {
	file, err := fh.Open()
	if err != nil {
		return err
	}
	u.SetAvatar(avatarType, file)
	return file.Close()
}

func (u *User) AvatarBase64(urlEnc bool) (avatarType, data string) {
	if urlEnc {
		return u.AvatarType, base64.URLEncoding.EncodeToString(u.Avatar)
	}
	return u.AvatarType, base64.StdEncoding.EncodeToString(u.Avatar)
}

///////////////////////////////////////////////////

func RemoveUser(uname string, lock bool) error {
	if lock {
		DbGrp.Lock()
		defer DbGrp.Unlock()
	}
	prefixes := [][]byte{
		[]byte("T" + SEP + uname + SEP),
		[]byte("F" + SEP + uname + SEP),
	}
	for _, prefix := range prefixes {
		n, err := bh.DeleteFirstObjectDB[User](prefix)
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
	DbGrp.Lock()
	defer DbGrp.Unlock()

	if err := RemoveUser(u.UName, false); err != nil {
		return err
	}
	return bh.UpsertOneObjectDB(u)
}

func LoadUser(uname string, active bool) (*User, bool, error) {
	DbGrp.Lock()
	defer DbGrp.Unlock()

	prefix := []byte("T" + SEP + uname + SEP)
	if !active {
		prefix = []byte("F" + SEP + uname + SEP)
	}
	u, err := bh.GetFirstObjectDB[User](prefix, nil)
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
	DbGrp.Lock()
	defer DbGrp.Unlock()

	return bh.GetObjectsDB([]byte(""), filter)
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
