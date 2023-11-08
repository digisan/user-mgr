package signin

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/digisan/user-mgr/db"
	u "github.com/digisan/user-mgr/user"
)

func TestOfflineMonitor(t *testing.T) {

	InitDB("../server-example/data/user")
	defer CloseDB()

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
		Hail("a")
		time.Sleep(1 * time.Second)
		Hail("a")
		time.Sleep(1 * time.Second)
		Hail("a")
		time.Sleep(1 * time.Second)
	}()

	go func() {
		Hail("b")
		time.Sleep(1 * time.Second)
		Hail("c")
		time.Sleep(1 * time.Second)
	}()

	go func() {
		time.Sleep(30 * time.Second)
		Hail("a")
	}()

	time.Sleep(1 * time.Minute)
}
