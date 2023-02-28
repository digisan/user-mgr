package signin

import (
	"context"
	"time"

	lk "github.com/digisan/logkit"
	u "github.com/digisan/user-mgr/user"
)

var (
	offlineTimeout time.Duration = 20 * time.Second
)

// frequently invoke it at Front-End. Interval should be less than 1 minute
func Trail(uname string) error {
	lk.Log("%v Heartbeats:", uname)
	_, err := u.RefreshOnline(uname)
	return err
}

func SetOfflineTimeout(period time.Duration) {
	offlineTimeout = period
}

func MonitorOffline(ctx context.Context, offline chan<- string, fnOnGotOffline func(uname string) error) {
	const interval = 10 * time.Second
	if offlineTimeout <= interval {
		offlineTimeout = 3 * interval
	}
	go func(ctx context.Context) {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				users, err := u.OnlineUsers()
				lk.WarnOnErr("%v", err)
				for _, usr := range users {
					usr, err := u.GetOnline(usr.Uname)
					lk.WarnOnErr("%v", err)
					if time.Since(usr.Tm) > offlineTimeout {
						if fnOnGotOffline != nil {
							lk.WarnOnErr("%v", fnOnGotOffline(usr.Uname))
						}
						offline <- usr.Uname
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}
