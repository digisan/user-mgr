package signin

import (
	"context"
	"fmt"
	"testing"
	"time"

	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/udb"
)

func TestInactiveMonitor(t *testing.T) {

	udb.OpenUserStorage("../data/user")
	defer udb.CloseUserStorage()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	removed := make(chan string, 1024)
	MonitorInactive(ctx, removed, 20*time.Second)
	go func() {
		for rm := range removed {
			fmt.Println("offline:", rm)
			lk.WarnOnErr("%v", udb.UserDB.RemoveOnlineUser(rm))
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
