package signin

import (
	"context"
	"time"

	lk "github.com/digisan/logkit"
	u "github.com/digisan/user-mgr/user"
)

const (
	CheckInterval = 10 * time.Second
)

func Trail(uname string) error {
	lk.Log("%v Hearbeats", uname)
	_, err := u.RefreshOnline(uname)
	return err
}

func MonitorInactive(ctx context.Context, inactive chan<- string, offlineTimeout time.Duration, fnOnGotInactive func(uname string) error) {

	if offlineTimeout <= CheckInterval {
		offlineTimeout = 2 * CheckInterval
	}

	go func(ctx context.Context) {
		ticker := time.NewTicker(CheckInterval)
		for {
			select {
			case <-ticker.C:
				users, err := u.OnlineUsers()
				lk.WarnOnErr("%v", err)
				for _, usr := range users {
					usr, err := u.GetOnline(usr.Uname)
					lk.WarnOnErr("%v", err)
					if time.Since(usr.Tm) > offlineTimeout {
						if fnOnGotInactive != nil {
							lk.WarnOnErr("%v", fnOnGotInactive(usr.Uname))
						}
						inactive <- usr.Uname
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}
