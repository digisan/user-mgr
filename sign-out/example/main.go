package main

import (
	"context"
	"fmt"
	"log"
	"time"

	. "github.com/digisan/user-mgr/db"
	si "github.com/digisan/user-mgr/sign-in"
	so "github.com/digisan/user-mgr/sign-out"
)

func main() {

	InitDB("../../server-example/data/user")
	defer CloseDB()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	offlineUsers := make(chan string, 1024)
	si.SetOfflineTimeout(20 * time.Second)
	si.MonitorOffline(ctx, offlineUsers, func(uname string) error { return so.Logout(uname) })
	go func() {
		for rm := range offlineUsers {
			fmt.Println("offline:", rm)
			if e := so.Logout(rm); e != nil {
				log.Fatalf("offline error @%s on %v", rm, e)
			}
		}
	}()

	time.Sleep(30 * time.Second)
}
