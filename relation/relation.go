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
		lk.FailOnErr("invalid flag [%d], only accept [enum.FOLLOWING enum.FOLLOWER enum.BLOCK enum.MUTED]", flag)
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
			lk.FailOnErr("invalid dbKey storage flag [%s], only accept ['FI' 'FR' 'B' 'M']", unames[0])
		}
	}
	panic("fatal error in dbKey for user relation")
}

var (
	mtx sync.Mutex
)

func RelAction(doFlag int, me, him string) (err error) {
	mtx.Lock()
	defer mtx.Unlock()

	if RelDB != nil {
		switch doFlag {
		case DO_FOLLOW:
			relMe, relHim := RelMgr(FOLLOWING, me), RelMgr(FOLLOWER, him)
			if str.NotIn(him, relMe.following...) {
				relMe.following = append(relMe.following, him)
				relHim.follower = append(relHim.follower, me)
				if err = RelDB.UpdateRel(FOLLOWING, relMe); err != nil {
					return
				}
				if err = RelDB.UpdateRel(FOLLOWER, relHim); err != nil {
					return
				}
			}

		case DO_UNFOLLOW:
			relMe, relHim := RelMgr(FOLLOWING, me), RelMgr(FOLLOWER, him)
			if str.In(him, relMe.following...) {
				str.DelOneEle(&relMe.following, him)
				str.DelOneEle(&relHim.follower, me)
				if err = RelDB.UpdateRel(FOLLOWING, relMe); err != nil {
					return
				}
				if err = RelDB.UpdateRel(FOLLOWER, relHim); err != nil {
					return
				}
			}

		case DO_BLOCK:
			relMe := RelMgr(BLOCKED, me)
			if str.NotIn(him, relMe.blocked...) {
				relMe.blocked = append(relMe.blocked, him)
				if err = RelDB.UpdateRel(FOLLOWING, relMe); err != nil {
					return
				}
			}

		case DO_UNBLOCK:
			relMe := RelMgr(BLOCKED, me)
			if str.In(him, relMe.blocked...) {
				str.DelOneEle(&relMe.blocked, him)
				if err = RelDB.UpdateRel(FOLLOWING, relMe); err != nil {
					return
				}
			}

		case DO_MUTE:
			relMe := RelMgr(MUTED, me)
			if str.NotIn(him, relMe.muted...) {
				relMe.muted = append(relMe.muted, him)
				if err = RelDB.UpdateRel(FOLLOWING, relMe); err != nil {
					return
				}
			}

		case DO_UNMUTE:
			relMe := RelMgr(MUTED, me)
			if str.In(him, relMe.muted...) {
				str.DelOneEle(&relMe.muted, him)
				if err = RelDB.UpdateRel(FOLLOWING, relMe); err != nil {
					return
				}
			}

		default:
			panic("invalid doFlag")
		}
	}
	return fmt.Errorf("RelDB is nil")
}
