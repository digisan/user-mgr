package signin

import (
	"fmt"
	"strings"
	"sync"
	"time"

	u "github.com/digisan/user-mgr/user"
	ur "github.com/digisan/user-mgr/user/registered"
)

var (
	smAccess = &sync.Map{}
	fnTS     = func(n int64) int64 {
		return time.Now().Unix() / n
	}
	smFrequently = &sync.Map{}
	smWaiting    = &sync.Map{}
)

func recordAccess(uname string, span int) (string, int) {
	key := fmt.Sprintf("%s-%d", uname, fnTS(int64(span)))
	if N, ok := smAccess.Load(key); ok {
		n := N.(int) + 1
		smAccess.Store(key, n)
		return key, n
	}
	n := 1
	smAccess.Store(key, n)
	return key, n
}

func delAccessRecord(uname string) int {
	keys := []string{}
	prefix := uname + "-"
	smAccess.Range(func(key, value any) bool {
		if k := key.(string); strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
		return true
	})
	for _, key := range keys {
		smAccess.Delete(key)
	}
	return len(keys)
}

func CheckFrequentlyAccess(uname string, spanSeconds, accessThreshold int) {
	if _, n := recordAccess(uname, spanSeconds); n >= accessThreshold {
		smFrequently.Store(uname, struct{}{})
		return
	}
}

func IsFrequentlyAccess(uname string) bool {
	_, ok := smFrequently.Load(uname)
	return ok
}

func RemoveFrequentlyAccessRecord(uname string, delay time.Duration) {
	go func() {
		if _, ok := smWaiting.Load(uname); ok {
			return
		}
		smWaiting.Store(uname, struct{}{})
		time.Sleep(delay)
		smFrequently.Delete(uname)
		delAccessRecord(uname)
		smWaiting.Delete(uname)
	}()
}

//////////////////////////////////////////////////////////////////////

// if return nil, which means user exists normally
func UserStatusIssue(login *ur.User) error {
	uname, email := login.UName, login.Email
	if u.UserExists(uname, email, true) {
		return nil
	}
	if u.UserExists(uname, email, false) {
		return fmt.Errorf("[%v] is dormant, cannot login", uname)
	}
	return fmt.Errorf("[%v] is not existing", uname)
}

// if successful, then update login user
func PwdOK(login *ur.User) bool {
	mPropVal := map[string]string{
		"uname": login.UName,
		"email": login.Email,
		"phone": login.Phone,
	}
	for prop, val := range mPropVal {
		if len(val) == 0 {
			continue
		}
		user, ok, err := u.LoadUserByUniProp(prop, val, true)
		if err == nil && ok && user.Password == login.Password {
			*login = *user
			return true
		}
	}
	return false
}
