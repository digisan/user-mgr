package relation

import (
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v3"
	bh "github.com/digisan/db-helper/badger-helper"
	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	u "github.com/digisan/user-mgr/user"
)

const (
	SEP = "^"
)

const (
	FOLLOWING int = iota
	FOLLOWER
	BLOCKED
	MUTED
)

const (
	DO_FOLLOW int = iota
	DO_UNFOLLOW
	DO_BLOCK
	DO_UNBLOCK
	DO_MUTE
	DO_UNMUTE
)

type Rel struct {
	uname       string
	mWithOthers map[int][]string
}

func (r Rel) String() string {
	return fmt.Sprintf("me: %v\n", r.uname) +
		fmt.Sprintf("following: %v\n", r.mWithOthers[FOLLOWING]) +
		fmt.Sprintf("follower: %v\n", r.mWithOthers[FOLLOWER]) +
		fmt.Sprintf("blocked: %v\n", r.mWithOthers[BLOCKED]) +
		fmt.Sprintf("muted: %v\n", r.mWithOthers[MUTED])
}

func (r *Rel) BadgerDB() *badger.DB {
	return DbGrp.Rel
}

func (r *Rel) Key() []byte {
	panic("Should NOT be used!")
	return nil
}

// at : [FOLLOWING FOLLOWER BLOCK MUTED]
func (r *Rel) Marshal(at any) (forKey, forValue []byte) {
	switch at {
	case FOLLOWING:
		forKey = []byte("FI" + SEP + r.uname)
		forValue = []byte(strings.Join(r.mWithOthers[FOLLOWING], SEP))
	case FOLLOWER:
		forKey = []byte("FR" + SEP + r.uname)
		forValue = []byte(strings.Join(r.mWithOthers[FOLLOWER], SEP))
	case BLOCKED:
		forKey = []byte("B" + SEP + r.uname)
		forValue = []byte(strings.Join(r.mWithOthers[BLOCKED], SEP))
	case MUTED:
		forKey = []byte("M" + SEP + r.uname)
		forValue = []byte(strings.Join(r.mWithOthers[MUTED], SEP))
	default:
		lk.FailOnErr("invalid 'at' [%d], only accept [FOLLOWING FOLLOWER BLOCK MUTED]", at)
	}
	return
}

func (r *Rel) Unmarshal(dbKey, dbVal []byte) (any, error) {
	unames := strings.Split(string(dbKey), SEP)
	others := strings.Split(string(dbVal), SEP)
	if len(unames) > 1 {
		r.uname = unames[1]
		switch unames[0] {
		case "FI":
			r.mWithOthers[FOLLOWING] = others
		case "FR":
			r.mWithOthers[FOLLOWER] = others
		case "B":
			r.mWithOthers[BLOCKED] = others
		case "M":
			r.mWithOthers[MUTED] = others
		default:
			lk.FailOnErr("invalid dbKey storage flag [%s], MUST be ['FI' 'FR' 'B' 'M']", unames[0])
		}
		return r, nil
	}
	panic("fatal error in dbKey for user relation")
}

////////////////////////////////////////////////////////////////////

func (r *Rel) HasFollowing(uname string) bool {
	return In(uname, r.mWithOthers[FOLLOWING]...)
}

func (r *Rel) HasFollower(uname string) bool {
	return In(uname, r.mWithOthers[FOLLOWER]...)
}

func (r *Rel) HasBlocked(uname string) bool {
	return In(uname, r.mWithOthers[BLOCKED]...)
}

func (r *Rel) HasMuted(uname string) bool {
	return In(uname, r.mWithOthers[MUTED]...)
}

////////////////////////////////////////////////////////////////////

// flag: [FOLLOWING FOLLOWER BLOCK MUTED]
func RemoveRel(uname string, flag int, lock bool) error {
	if lock {
		DbGrp.Lock()
		defer DbGrp.Unlock()
	}

	mKey := map[int][]byte{
		FOLLOWING: []byte("FI" + SEP + uname),
		FOLLOWER:  []byte("FR" + SEP + uname),
		BLOCKED:   []byte("B" + SEP + uname),
		MUTED:     []byte("M" + SEP + uname),
	}

	key, ok := mKey[flag]
	lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("invalid flag"))

	_, err := bh.DeleteOneObjectDB[Rel](key)
	return err
}

// flag: [FOLLOWING FOLLOWER BLOCK MUTED]
func UpdateRel(r *Rel, flag int) error {
	DbGrp.Lock()
	defer DbGrp.Unlock()

	if err := RemoveRel(r.uname, flag, false); err != nil {
		return err
	}
	return bh.UpsertPartObjectDB(r, flag)
}

// flag: [FOLLOWING FOLLOWER BLOCK MUTED]
func LoadRel(uname string, flag int) (*Rel, bool, error) {
	DbGrp.Lock()
	defer DbGrp.Unlock()

	mKey := map[int][]byte{
		FOLLOWING: []byte("FI" + SEP + uname),
		FOLLOWER:  []byte("FR" + SEP + uname),
		BLOCKED:   []byte("B" + SEP + uname),
		MUTED:     []byte("M" + SEP + uname),
	}

	key, ok := mKey[flag]
	lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("invalid flag, only accept [FOLLOWING FOLLOWER BLOCKED MUTED]"))

	r, err := bh.GetOneObjectDB[Rel](key)
	return r, r != nil && r.uname != "", err
}

// flag: [FOLLOWING FOLLOWER BLOCK MUTED]
func ListRel(uname string, flag int) []string {
	if r, ok, err := LoadRel(uname, flag); err == nil && ok {
		return r.mWithOthers[flag]
	}
	return nil
}

func RelMgr(uname string, flags ...int) *Rel {
	r := &Rel{uname: uname}
	for _, flag := range flags {
		r.mWithOthers[flag] = ListRel(uname, flag)
	}
	return r
}

func relAction(me string, doFlag int, whom string, lock bool) error {
	if lock {
		DbGrp.Lock()
		defer DbGrp.Unlock()
	}

	if DbGrp.Rel != nil && u.DbGrp != nil && u.DbGrp.Reg != nil {

		if u.UserExists(whom, "", false) {

			switch doFlag {
			case DO_FOLLOW, DO_UNFOLLOW:

				meFollowing := ListRel(me, FOLLOWING)
				whomFollower := ListRel(whom, FOLLOWER)
				did := false

				if doFlag == DO_FOLLOW && NotIn(whom, meFollowing...) {
					meFollowing = append(meFollowing, whom)
					whomFollower = append(whomFollower, me)
					did = true
				} else if doFlag == DO_UNFOLLOW && In(whom, meFollowing...) {
					DelOneEle(&meFollowing, whom)
					DelOneEle(&whomFollower, me)
					did = true
				}
				if did {
					m := map[int][]string{
						FOLLOWING: meFollowing,
						FOLLOWER:  whomFollower,
					}
					if err := UpdateRel(&Rel{me, m}, FOLLOWING); err != nil {
						return err
					}
					if err := UpdateRel(&Rel{whom, m}, FOLLOWER); err != nil {
						return err
					}
				}

			case DO_BLOCK, DO_UNBLOCK:

				blocked := ListRel(me, BLOCKED)
				did := false

				if doFlag == DO_BLOCK && NotIn(whom, blocked...) {
					blocked = append(blocked, whom)
					// other actions
					{
						relAction(me, DO_UNFOLLOW, whom, false)
					}
					//
					did = true
				} else if doFlag == DO_UNBLOCK && In(whom, blocked...) {
					DelOneEle(&blocked, whom)
					did = true
				}
				if did {
					m := map[int][]string{
						BLOCKED: blocked,
					}
					return UpdateRel(&Rel{me, m}, BLOCKED)
				}

			case DO_MUTE, DO_UNMUTE:

				muted := ListRel(me, MUTED)
				did := false

				if doFlag == DO_MUTE && NotIn(whom, muted...) {
					muted = append(muted, whom)
					// other actions
					{

					}
					//
					did = true
				} else if doFlag == DO_UNMUTE && In(whom, muted...) {
					DelOneEle(&muted, whom)
					did = true
				}
				if did {
					m := map[int][]string{
						MUTED: muted,
					}
					return UpdateRel(&Rel{me, m}, MUTED)
				}

			default:
				panic("invalid doFlag, only accept [DO_FOLLOW DO_UNFOLLOW DO_BLOCK DO_UNBLOCK DO_MUTE DO_UNMUTE]")
			}

		} else {
			return fmt.Errorf("%s is not registered", whom)
		}

	} else {
		return fmt.Errorf("DbGrp.Rel or DbGrp.Reg is nil")
	}
	return nil
}

// doFlag: [DO_FOLLOW, DO_UNFOLLOW, DO_BLOCK, DO_UNBLOCK, DO_MUTE, DO_UNMUTE]
func RelAction(me string, doFlag int, whom string) error {
	return relAction(me, doFlag, whom, true)
}
