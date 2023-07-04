package registered

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/crypto"
	lk "github.com/digisan/logkit"
	. "github.com/digisan/user-mgr/cst"
	"github.com/digisan/user-mgr/db"
	. "github.com/digisan/user-mgr/util"
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
	KO_RegTime
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
		KO_RegTime:   &u.RegTime,
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
	return db.DbGrp.Registered
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
			switch ip {
			case 0:
				segs = bytes.Split(param.in, []byte(SEP))
				u.key = *(*[16]byte)(segs[KO_Key])
			case 1:
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

// DO NOT  apply below functions directly to user from claims !!!
// Apply them to user from db or cache !!!

func (u *User) IsActive() bool {
	return u.Active
}

func (u *User) StampRegTime() {
	u.RegTime = time.Now()
}

func (u *User) SinceJoined() time.Duration {
	t := &time.Time{}
	t.UnmarshalText([]byte(u.RegTime.Format(time.RFC3339)))
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

func (u *User) SetAvatar(r io.Reader, avatarType string) {
	u.AvatarType = avatarType
	u.Avatar = StreamToBytes(r)
}

// 'avatarType' --- e.g. image/png, 'fh' --- FormFile('param')
// for example '<img src="data:image/png;base64,******/>'
func (u *User) SetAvatarByFormFile(fh *multipart.FileHeader, x, y, w, h int) error {
	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// prepare a temp file for writing, then delete temp file
	tempAvatar := fmt.Sprintf(`/tmp/temp-avatar-%v`, u.UName)
	writer, err := os.Create(tempAvatar)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempAvatar)
	defer writer.Close()
	// write avatar to temp file from file header
	if _, err = io.Copy(writer, file); err != nil {
		return err
	}

	// crop temp image and set it as avatar, then delete cropped file
	if op, err := CropImage(tempAvatar, fmt.Sprintf(`crop:%d,%d,%d,%d`, x, y, w, h), "png"); err == nil && len(op) != 0 {
		cropped, err := os.Open(op)
		if err != nil {
			return err
		}
		defer os.RemoveAll(op)
		defer cropped.Close()

		avatarType := "image/" + strings.TrimPrefix(filepath.Ext(fh.Filename), ".")
		u.SetAvatar(cropped, avatarType)
	}
	return nil
}

func (u *User) AvatarBase64(urlEnc bool) (data, avatarType string) {
	if urlEnc {
		return base64.URLEncoding.EncodeToString(u.Avatar), u.AvatarType
	}
	return base64.StdEncoding.EncodeToString(u.Avatar), u.AvatarType
}
