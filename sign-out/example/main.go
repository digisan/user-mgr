package main

import (
	"context"
	"fmt"
	"time"

	si "github.com/digisan/user-mgr/sign-in"
	so "github.com/digisan/user-mgr/sign-out"
	u "github.com/digisan/user-mgr/user"
)

func main() {

	u.InitDB("../../data/user")
	defer u.CloseDB()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inactiveUsers := make(chan string, 1024)
	si.MonitorInactive(ctx, inactiveUsers, 20*time.Second, func(uname string) error { return so.Logout("QMiao") })
	go func() {
		for rm := range inactiveUsers {
			fmt.Println("offline:", rm)
		}
	}()

	time.Sleep(60 * time.Second)
}
