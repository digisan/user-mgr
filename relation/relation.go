package relation

import (
	"fmt"
	"strings"
	"sync"

	"github.com/digisan/go-generics/str"
	lk "github.com/digisan/logkit"
	. "github.com/digisan/user-mgr/relation/enum"
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
	return str.In(uname, r.following...)
}

func (r *Rel) HasFollower(uname string) bool {
	return str.In(uname, r.follower...)
}

func (r *Rel) HasBlocked(uname string) bool {
	return str.In(uname, r.blocked...)
}

func (r *Rel) HasMuted(uname string) bool {
	return str.In(uname, r.muted...)
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

func relAction(me string, doFlag int, him string, lock bool) (err error) {
	if lock {
		mtx.Lock()
		defer mtx.Unlock()
	}

	if RelDB != nil {

		switch doFlag {
		case DO_FOLLOW, DO_UNFOLLOW:

			meFollowing := RelContent(me, FOLLOWING)
			himFollower := RelContent(him, FOLLOWER)
			did := false

			if doFlag == DO_FOLLOW && str.NotIn(him, meFollowing...) {
				meFollowing = append(meFollowing, him)
				himFollower = append(himFollower, me)
				did = true
			} else if doFlag == DO_UNFOLLOW && str.In(him, meFollowing...) {
				str.DelOneEle(&meFollowing, him)
				str.DelOneEle(&himFollower, me)
				did = true
			}
			if did {
				if err = RelDB.UpdateRel(FOLLOWING, &Rel{uname: me, following: meFollowing}); err != nil {
					return
				}
				if err = RelDB.UpdateRel(FOLLOWER, &Rel{uname: him, follower: himFollower}); err != nil {
					return
				}
			}

		case DO_BLOCK, DO_UNBLOCK:

			blocked := RelContent(me, BLOCKED)
			did := false

			if doFlag == DO_BLOCK && str.NotIn(him, blocked...) {
				blocked = append(blocked, him)
				// other actions
				{
					relAction(me, DO_UNFOLLOW, him, false)
				}
				//
				did = true
			} else if doFlag == DO_UNBLOCK && str.In(him, blocked...) {
				str.DelOneEle(&blocked, him)
				did = true
			}
			if did {
				return RelDB.UpdateRel(BLOCKED, &Rel{uname: me, blocked: blocked})
			}

		case DO_MUTE, DO_UNMUTE:

			muted := RelContent(me, MUTED)
			did := false

			if doFlag == DO_MUTE && str.NotIn(him, muted...) {
				muted = append(muted, him)
				// other actions
				{

				}
				//
				did = true
			} else if doFlag == DO_UNMUTE && str.In(him, muted...) {
				str.DelOneEle(&muted, him)
				did = true
			}
			if did {
				return RelDB.UpdateRel(MUTED, &Rel{uname: me, muted: muted})
			}

		default:
			panic("invalid doFlag, only accept [DO_FOLLOW DO_UNFOLLOW DO_BLOCK DO_UNBLOCK DO_MUTE DO_UNMUTE]")
		}

	} else {
		return fmt.Errorf("RelDB is nil")
	}

	return
}

func RelAction(me string, doFlag int, him string) (err error) {
	return relAction(me, doFlag, him, true)
}
