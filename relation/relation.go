package relation

import (
	"fmt"
	"strings"
	"sync"

	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	. "github.com/digisan/user-mgr/relation/enum"
	"github.com/digisan/user-mgr/udb"
)

const (
	SEP = "^"
)

type Rel struct {
	uname     string
	following []string
	follower  []string
	blocked   []string
	muted     []string
}

func (r Rel) String() string {
	return fmt.Sprintf("me: %v\n", r.uname) +
		fmt.Sprintf("following: %v\n", r.following) +
		fmt.Sprintf("follower: %v\n", r.follower) +
		fmt.Sprintf("blocked: %v\n", r.blocked) +
		fmt.Sprintf("muted: %v\n", r.muted)
}

func (r *Rel) MarshalTo(flag int) (forKey, forValue []byte) {
	switch flag {
	case FOLLOWING:
		forKey = []byte("FI" + SEP + r.uname)
		forValue = []byte(strings.Join(r.following, SEP))
	case FOLLOWER:
		forKey = []byte("FR" + SEP + r.uname)
		forValue = []byte(strings.Join(r.follower, SEP))
	case BLOCKED:
		forKey = []byte("B" + SEP + r.uname)
		forValue = []byte(strings.Join(r.blocked, SEP))
	case MUTED:
		forKey = []byte("M" + SEP + r.uname)
		forValue = []byte(strings.Join(r.muted, SEP))
	default:
		lk.FailOnErr("invalid flag [%d], only accept [FOLLOWING FOLLOWER BLOCK MUTED]", flag)
	}
	return
}

func (r *Rel) UnmarshalFrom(dbKey, dbVal []byte) int {
	unames := strings.Split(string(dbKey), SEP)
	others := strings.Split(string(dbVal), SEP)
	if len(unames) > 1 {
		r.uname = unames[1]
		switch unames[0] {
		case "FI":
			r.following = others
			return FOLLOWING
		case "FR":
			r.follower = others
			return FOLLOWER
		case "B":
			r.blocked = others
			return BLOCKED
		case "M":
			r.muted = others
			return MUTED
		default:
			lk.FailOnErr("invalid dbKey storage flag [%s], MUST be ['FI' 'FR' 'B' 'M']", unames[0])
		}
	}
	panic("fatal error in dbKey for user relation")
}

var (
	mtx sync.Mutex
)

func (r *Rel) HasFollowing(uname string) bool {
	return In(uname, r.following...)
}

func (r *Rel) HasFollower(uname string) bool {
	return In(uname, r.follower...)
}

func (r *Rel) HasBlocked(uname string) bool {
	return In(uname, r.blocked...)
}

func (r *Rel) HasMuted(uname string) bool {
	return In(uname, r.muted...)
}

func RelMgr(uname string, flags ...int) *Rel {
	rel := &Rel{uname: uname}
	m := map[int]*[]string{
		FOLLOWING: &rel.following,
		FOLLOWER:  &rel.follower,
		BLOCKED:   &rel.blocked,
		MUTED:     &rel.muted,
	}
	if len(flags) == 0 {
		flags = []int{FOLLOWING, FOLLOWER, BLOCKED, MUTED}
	}
	for _, flag := range flags {
		ptr, ok := m[flag]
		lk.FailOnErrWhen(!ok, "%v", fmt.Errorf("invalid flag"))
		*ptr = RelContent(uname, flag)
	}
	return rel
}

func relAction(me string, doFlag int, whom string, lock bool) (err error) {
	if lock {
		mtx.Lock()
		defer mtx.Unlock()
	}

	if RelDB != nil && udb.UserDB != nil {

		if udb.UserDB.UserExists(whom, "", false) {

			switch doFlag {
			case DO_FOLLOW, DO_UNFOLLOW:

				meFollowing := RelContent(me, FOLLOWING)
				whomFollower := RelContent(whom, FOLLOWER)
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
					if err = RelDB.UpdateRel(FOLLOWING, &Rel{uname: me, following: meFollowing}); err != nil {
						return
					}
					if err = RelDB.UpdateRel(FOLLOWER, &Rel{uname: whom, follower: whomFollower}); err != nil {
						return
					}
				}

			case DO_BLOCK, DO_UNBLOCK:

				blocked := RelContent(me, BLOCKED)
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
					return RelDB.UpdateRel(BLOCKED, &Rel{uname: me, blocked: blocked})
				}

			case DO_MUTE, DO_UNMUTE:

				muted := RelContent(me, MUTED)
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
					return RelDB.UpdateRel(MUTED, &Rel{uname: me, muted: muted})
				}

			default:
				panic("invalid doFlag, only accept [DO_FOLLOW DO_UNFOLLOW DO_BLOCK DO_UNBLOCK DO_MUTE DO_UNMUTE]")
			}

		} else {
			return fmt.Errorf("%s is not registered", whom)
		}

	} else {
		return fmt.Errorf("RelDB or UserDB is nil")
	}

	return
}

func RelAction(me string, doFlag int, whom string) (err error) {
	return relAction(me, doFlag, whom, true)
}
