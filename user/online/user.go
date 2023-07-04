package online

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	lk "github.com/digisan/logkit"
	. "github.com/digisan/user-mgr/cst"
	"github.com/digisan/user-mgr/db"
)

// if modified, change 1. KO_OL_***, 2. mFldAddr, 3. 'auto-tags.go', 4. 'validator.go' in sign-up.
type User struct {
	// key
	Uname string
	// value
	Tm time.Time
}

func NewUser(uname string) *User {
	return &User{uname, time.Now().UTC()}
}

func (u User) String() string {
	return fmt.Sprintf("%s @ %v\n", u.Uname, u.Tm)
}

// db key order
const (
	KO_UName int = iota
	KO_END
)

// db value order
const (
	VO_Tm int = iota
	VO_END
)

func (u *User) KeyFieldAddr(ko int) any {
	mFldAddr := map[int]any{
		KO_UName: &u.Uname,
	}
	return mFldAddr[ko]
}

func (u *User) ValFieldAddr(vo int) any {
	mFldAddr := map[int]any{
		VO_Tm: &u.Tm,
	}
	return mFldAddr[vo]
}

////////////////////////////////////////////////////

func (u *User) BadgerDB() *badger.DB {
	return db.DbGrp.Online
}

func (u *User) Key() []byte {
	var (
		sb = &strings.Builder{}
	)
	for i := 0; i < KO_END; i++ {
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

func (u *User) Value() []byte {
	var (
		sb = &strings.Builder{}
	)
	for i := 0; i < VO_END; i++ {
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
			for i, seg := range bytes.Split(param.in, []byte(SEP)) {
				if (ip == 0 && i == KO_END) || (ip == 1 && i == VO_END) {
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
