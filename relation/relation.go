package relation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v4"
	bh "github.com/digisan/db-helper/badger"
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
	FOLLOW int = iota
	UNFOLLOW
	BLOCK
	UNBLOCK
	MUTE
	UNMUTE
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
	if len(dbVal) == 0 {
		others = []string{}
	}
	if len(unames) > 1 {
		r.uname = unames[1]
		r.mWithOthers = make(map[int][]string)
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

	_, err := bh.DeleteOneObject[Rel](key)
	return err
}

// flag: [FOLLOWING FOLLOWER BLOCK MUTED]
func UpdateRel(r *Rel, flag int, lock bool) error {
	if lock {
		DbGrp.Lock()
		defer DbGrp.Unlock()
	}

	if err := RemoveRel(r.uname, flag, false); err != nil {
		return err
	}
	return bh.UpsertPartObject(r, flag)
}

// flag: [FOLLOWING FOLLOWER BLOCK MUTED]
func LoadRel(uname string, flag int, lock bool) (*Rel, bool, error) {
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
	lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("invalid flag, only accept [FOLLOWING FOLLOWER BLOCKED MUTED]"))

	r, err := bh.GetOneObject[Rel](key)
	return r, r != nil && r.uname != "", err
}

// flag: [FOLLOWING FOLLOWER BLOCK MUTED]
func ListRel(uname string, flag int, lock bool) ([]string, error) {
	r, ok, err := LoadRel(uname, flag, lock)
	if err != nil {
		return nil, err
	}
	if ok {
		return r.mWithOthers[flag], nil
	}
	return []string{}, nil
}

func RelMgr(uname string, flags ...int) (*Rel, error) {
	r := &Rel{uname: uname, mWithOthers: make(map[int][]string)}
	if len(flags) == 0 {
		flags = []int{FOLLOWING, FOLLOWER, BLOCKED, MUTED}
	}
	for _, flag := range flags {
		names, err := ListRel(uname, flag, true)
		if err != nil {
			return nil, err
		}
		r.mWithOthers[flag] = names
	}
	return r, nil
}

func relAction(me string, flag int, whom string, lock bool) error {
	if lock {
		DbGrp.Lock()
		defer DbGrp.Unlock()
	}

	if len(whom) == 0 {
		return nil
	}

	if DbGrp.Rel != nil && u.DbGrp != nil && u.DbGrp.Registered != nil {
		if !u.UserExists(me, "", false) {
			return fmt.Errorf("'%s' is not registered", me)
		}
		if !u.UserExists(whom, "", false) {
			return fmt.Errorf("'%s' is not registered", whom)
		}

		switch flag {
		case FOLLOW, UNFOLLOW:
			meFollowing, err := ListRel(me, FOLLOWING, false)
			if err != nil {
				return err
			}
			whomFollower, err := ListRel(whom, FOLLOWER, false)
			if err != nil {
				return err
			}
			did := false

			if flag == FOLLOW && NotIn(whom, meFollowing...) {
				meFollowing = append(meFollowing, whom)
				whomFollower = append(whomFollower, me)
				did = true
			} else if flag == UNFOLLOW && In(whom, meFollowing...) {
				DelOneEle(&meFollowing, whom)
				DelOneEle(&whomFollower, me)
				did = true
			}
			if did {
				m := map[int][]string{
					FOLLOWING: meFollowing,
					FOLLOWER:  whomFollower,
				}
				if err := UpdateRel(&Rel{me, m}, FOLLOWING, false); err != nil {
					return err
				}
				if err := UpdateRel(&Rel{whom, m}, FOLLOWER, false); err != nil {
					return err
				}
			}

		case BLOCK, UNBLOCK:
			blocked, err := ListRel(me, BLOCKED, lock)
			if err != nil {
				return err
			}
			did := false

			if flag == BLOCK && NotIn(whom, blocked...) {
				blocked = append(blocked, whom)
				// other actions
				{
					relAction(me, UNFOLLOW, whom, false)
				}
				//
				did = true
			} else if flag == UNBLOCK && In(whom, blocked...) {
				DelOneEle(&blocked, whom)
				did = true
			}
			if did {
				m := map[int][]string{
					BLOCKED: blocked,
				}
				return UpdateRel(&Rel{me, m}, BLOCKED, false)
			}

		case MUTE, UNMUTE:
			muted, err := ListRel(me, MUTED, lock)
			if err != nil {
				return err
			}
			did := false

			if flag == MUTE && NotIn(whom, muted...) {
				muted = append(muted, whom)
				// other actions
				{

				}
				//
				did = true
			} else if flag == UNMUTE && In(whom, muted...) {
				DelOneEle(&muted, whom)
				did = true
			}
			if did {
				m := map[int][]string{
					MUTED: muted,
				}
				return UpdateRel(&Rel{me, m}, MUTED, false)
			}

		default:
			panic("invalid flag, only accept [FOLLOW UNFOLLOW BLOCK UNBLOCK MUTE UNMUTE]")
		}

	} else {
		return fmt.Errorf("DbGrp.Rel or DbGrp.Reg is nil")
	}
	return nil
}

// flag: [FOLLOW, UNFOLLOW, BLOCK, UNBLOCK, MUTE, UNMUTE]
// whom: if "ALL", clear all add flag [UNFOLLOW UNBLOCK UNMUTE]
func RelAction(me string, flag int, whom string) error {
	if whom == "ALL" {
		m := map[int]int{
			UNFOLLOW: FOLLOWING,
			UNBLOCK:  BLOCKED,
			UNMUTE:   MUTED,
		}
		status, ok := m[flag]
		if !ok {
			return errors.New("invalid flag when 'whom == ALL', only accept [UNFOLLOW UNBLOCK UNMUTE]")
		}
		names, err := ListRel(me, status, true)
		if err != nil {
			return err
		}
		for _, name := range names {
			if err := relAction(me, flag, name, true); err != nil {
				return err
			}
		}
		return nil
	}
	return relAction(me, flag, whom, true)
}
