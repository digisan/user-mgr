package signin

import (
	"context"
	"time"

	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/udb"
)

const (
	CheckInterval = 10 * time.Second
)

func Trail(uname string) error {
	return udb.UserDB.RefreshOnline(uname)
}

func MonitorInactive(ctx context.Context, inactive chan<- string, offlineTimeout time.Duration) {

	if offlineTimeout <= CheckInterval {
		offlineTimeout = 2 * CheckInterval
	}

	go func(ctx context.Context) {
		ticker := time.NewTicker(CheckInterval)
		for {
			select {
			case <-ticker.C:
				unames, err := udb.UserDB.OnlineUsers()
				lk.WarnOnErr("%v", err)
				for _, uname := range unames {
					lastTm, err := udb.UserDB.GetOnline(uname)
					lk.WarnOnErr("%v", err)
					if time.Since(lastTm) > offlineTimeout {
						// lk.WarnOnErr("%v", udb.UserDB.RemoveOnlineUser(uname)) // *** let invoker to do more
						inactive <- uname
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}
