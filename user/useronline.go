package user

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
	bh "github.com/digisan/db-helper/badger"
	lk "github.com/digisan/logkit"
)

// if modified, change 1. KO_OL_***, 2. mFldAddr, 3. 'auto-tags.go', 4. 'validator.go' in sign-up.
type UserOnline struct {
	// key
	Uname string
	// value
	Tm time.Time
}

func NewUserOnline(uname string) *UserOnline {
	return &UserOnline{uname, time.Now().UTC()}
}

func (u UserOnline) String() string {
	return fmt.Sprintf("%s @ %v\n", u.Uname, u.Tm)
}

// db key order
const (
	KO_OL_UName int = iota
	KO_OL_END
)

// db value order
const (
	VO_OL_Tm int = iota
	VO_OL_END
)

func (u *UserOnline) KeyFieldAddr(ko int) any {
	mFldAddr := map[int]any{
		KO_OL_UName: &u.Uname,
	}
	return mFldAddr[ko]
}

func (u *UserOnline) ValFieldAddr(vo int) any {
	mFldAddr := map[int]any{
		VO_OL_Tm: &u.Tm,
	}
	return mFldAddr[vo]
}

////////////////////////////////////////////////////

func (u *UserOnline) BadgerDB() *badger.DB {
	return DbGrp.Online
}

func (u *UserOnline) Key() []byte {
	var (
		sb = &strings.Builder{}
	)
	for i := 0; i < KO_OL_END; i++ {
		if i > 0 {
			sb.WriteString(SEP)
		}
		switch v := u.KeyFieldAddr(i).(type) {
		case *string:
			sb.WriteString(*v)
		default:
			panic("need more type for marshaling key")
		}
	}
	return []byte(sb.String())
}

func (u *UserOnline) Value() []byte {
	var (
		sb = &strings.Builder{}
	)
	for i := 0; i < VO_OL_END; i++ {
		if i > 0 {
			sb.WriteString(SEP)
		}
		switch v := u.ValFieldAddr(i).(type) {
		case *time.Time:
			encoding, err := (*v).MarshalBinary()
			// lk.Debug(" --------------- encoding len: %d", len(encoding))
			lk.FailOnErr("%v", err)
			sb.Write(encoding)
		default:
			panic("need more type for marshaling value")
		}
	}
	return []byte(sb.String())
}

func (u *UserOnline) Marshal(at any) (forKey, forValue []byte) {
	return u.Key(), u.Value()
}

func (u *UserOnline) Unmarshal(dbKey, dbVal []byte) (any, error) {
	params := []struct {
		in        []byte
		fnFldAddr func(int) any
	}{
		{dbKey, u.KeyFieldAddr},
		{dbVal, u.ValFieldAddr},
	}
	for ip, param := range params {
		if len(param.in) > 0 {
			for i, seg := range bytes.Split(param.in, []byte(SEP)) {
				if (ip == 0 && i == KO_OL_END) || (ip == 1 && i == VO_OL_END) {
					break
				}
				switch v := param.fnFldAddr(i).(type) {
				case *string:
					*v = string(seg)
				case *time.Time:
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

func GetOnline(uname string) (*UserOnline, error) {
	DbGrp.Lock()
	defer DbGrp.Unlock()

	return bh.GetOneObject[UserOnline]([]byte(uname))
}

func RefreshOnline(uname string) (*UserOnline, error) {
	DbGrp.Lock()
	defer DbGrp.Unlock()

	u := NewUserOnline(uname)
	return u, bh.UpsertOneObject(u)
}

func RmOnline(uname string) (int, error) {
	DbGrp.Lock()
	defer DbGrp.Unlock()

	return bh.DeleteOneObject[UserOnline]([]byte(uname))
}

func OnlineUsers() ([]*UserOnline, error) {
	DbGrp.Lock()
	defer DbGrp.Unlock()

	return bh.GetObjects[UserOnline]([]byte(""), nil)
}
