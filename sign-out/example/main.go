package main

import (
	"context"
	"fmt"
	"time"

	lk "github.com/digisan/logkit"
	si "github.com/digisan/user-mgr/sign-in"
	so "github.com/digisan/user-mgr/sign-out"
	"github.com/digisan/user-mgr/udb"
)

func main() {
	udb.OpenUserStorage("../../data/user")
	defer udb.CloseUserStorage()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inactiveUsers := make(chan string, 1024)
	si.MonitorInactive(ctx, inactiveUsers, 20*time.Second)
	go func() {
		for rm := range inactiveUsers {
			fmt.Println("offline:", rm)
			lk.WarnOnErr("%v", so.Logout("QMiao"))
		}
	}()

	time.Sleep(60 * time.Second)
}
