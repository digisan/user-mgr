package user

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/tool"
)

// if modified, change 1. MOK_***, 2. mFldAddr, 3. 'auto-tags.go', 4. 'validator.go' in sign-up.
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
	MOK_Active int = iota
	// Core
	MOK_UName
	MOK_Email
	MOK_Password
	MOK_Key
	// Profile
	MOK_Name
	MOK_Phone
	MOK_Country
	MOK_City
	MOK_Addr
	MOK_PersonalIDType
	MOK_PersonalID
	MOK_Gender
	MOK_DOB
	MOK_Position
	MOK_Title
	MOK_Employer
	MOK_Bio
	MOK_AvatarType
	// Admin
	MOK_Regtime
	MOK_SysRole
	MOK_MemLevel
	MOK_MemExpire
	MOK_Official
	MOK_Tags
	//
	MOK_END
)

// db value order
const (
	// Profile
	MOV_Avatar int = iota
	//
	MOV_END
)

func (u *User) KeyFieldAddr(mok int) any {
	mFldAddr := map[int]any{
		// Core
		MOK_UName:    &u.UName,
		MOK_Email:    &u.Email,
		MOK_Password: &u.Password,
		MOK_Key:      &u.key,
		// Profile
		MOK_Name:           &u.Name,
		MOK_Phone:          &u.Phone,
		MOK_Country:        &u.Country,
		MOK_City:           &u.City,
		MOK_Addr:           &u.Addr,
		MOK_PersonalIDType: &u.PersonalIDType,
		MOK_PersonalID:     &u.PersonalID,
		MOK_Gender:         &u.Gender,
		MOK_DOB:            &u.DOB,
		MOK_Position:       &u.Position,
		MOK_Title:          &u.Title,
		MOK_Employer:       &u.Employer,
		MOK_Bio:            &u.Bio,
		MOK_AvatarType:     &u.AvatarType,
		// Admin
		MOK_Active:    &u.Active,
		MOK_Regtime:   &u.Regtime,
		MOK_SysRole:   &u.SysRole,
		MOK_MemLevel:  &u.MemLevel,
		MOK_MemExpire: &u.MemExpire,
		MOK_Official:  &u.Official,
		MOK_Tags:      &u.Tags,
	}
	return mFldAddr[mok]
}

func (u *User) ValFieldAddr(mov int) any {
	mFldAddr := map[int]any{
		// Profile
		MOV_Avatar: &u.Avatar,
	}
	return mFldAddr[mov]
}

////////////////////////////////////////////////////

var secret = []int{
	MOK_Email,
	MOK_Password,
	MOK_Name,
	MOK_Phone,
	MOK_Country,
	MOK_City,
	MOK_Addr,
	MOK_PersonalIDType,
	MOK_PersonalID,
	MOK_Gender,
	MOK_DOB,
	MOK_Position,
	MOK_Title,
	MOK_Employer,
}

func (u *User) Marshal() (forKey, forValue []byte) {

	key := u.GenKey()

	params := []struct {
		end       int
		fnFldAddr func(int) any
		out       *[]byte
	}{
		{
			end:       MOK_END,
			fnFldAddr: u.KeyFieldAddr,
			out:       &forKey,
		},
		{
			end:       MOV_END,
			fnFldAddr: u.ValFieldAddr,
			out:       &forValue,
		},
	}
	for ip, param := range params {
		sb := &strings.Builder{}
		for i := 0; i < param.end; i++ {
			if i > 0 {
				sb.WriteString(SEP)
			}
			if ip == 0 && In(i, secret...) {
				sb.Write(tool.Encrypt((*param.fnFldAddr(i).(*string)), key[:]))
				continue
			}
			switch v := param.fnFldAddr(i).(type) {
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
			default:
				panic("Marshal Error Type")
			}
		}
		*param.out = []byte(sb.String())
	}
	return
}

func (u *User) Unmarshal(dbKey, dbVal []byte) {

	params := []struct {
		in        []byte
		fnFldAddr func(int) any
	}{
		{
			in:        dbKey,
			fnFldAddr: u.KeyFieldAddr,
		},
		{
			in:        dbVal,
			fnFldAddr: u.ValFieldAddr,
		},
	}
	for ip, param := range params {
		if len(param.in) > 0 {

			var segs [][]byte

			if ip == 0 {
				segs = bytes.Split(param.in, []byte(SEP))
				u.key = *(*[16]byte)(segs[MOK_Key])
			} else if ip == 1 {
				segs = [][]byte{param.in} // dbVal is one whole block
			}

			for i, seg := range segs {
				if (ip == 0 && i == MOK_END) || (ip == 1 && i == MOV_END) {
					break
				}
				if ip == 0 && In(i, secret...) {
					if u.key != [16]byte{} {
						*param.fnFldAddr(i).(*string) = tool.Decrypt(seg, u.key[:])
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
				default:
					panic("Unmarshal Error Type")
				}
			}
		}
	}
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
	u.Avatar = tool.StreamToByte(r)
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
