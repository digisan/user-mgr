package signin

import (
	"context"
	"fmt"
	"testing"
	"time"

	u "github.com/digisan/user-mgr/user"
)

func TestOfflineMonitor(t *testing.T) {

	u.InitDB("../server-example/data/user")
	defer u.CloseDB()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	removed := make(chan string, 1024)
	SetOfflineTimeout(20 * time.Second)
	MonitorOffline(ctx, removed, func(uname string) error {
		_, err := u.RmOnline(uname)
		return err
	})
	go func() {
		for rm := range removed {
			fmt.Println("offline:", rm)
		}
	}()

	go func() {
		Trail("a")
		time.Sleep(1 * time.Second)
		Trail("a")
		time.Sleep(1 * time.Second)
		Trail("a")
		time.Sleep(1 * time.Second)
	}()

	go func() {
		Trail("b")
		time.Sleep(1 * time.Second)
		Trail("c")
		time.Sleep(1 * time.Second)
	}()

	go func() {
		time.Sleep(30 * time.Second)
		Trail("a")
	}()

	time.Sleep(1 * time.Minute)
}
