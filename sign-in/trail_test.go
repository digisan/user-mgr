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

func TestTimeStamp(t *testing.T) {

	fmt.Println(recordAccess("qmiao", 3))
	time.Sleep(1 * time.Second)
	fmt.Println(recordAccess("qmiao", 3))
	time.Sleep(2 * time.Second)
	fmt.Println(recordAccess("qmiao", 4))

	fmt.Println("---------")

	delAccessRecord("qmiao")

	smAccess.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})
}

func TestLockCheck(t *testing.T) {
	const USER = "test"
	uname := USER

	for i := 0; i < 20; i++ {

		fmt.Printf("%s trying... %02d\n", uname, i)

		if IsFrequentlyAccess(uname) {
			delay := 5 * time.Second
			fmt.Printf("locked, wait %v\n", delay)
			RemoveFrequentlyAccessRecord(uname, delay)
			time.Sleep(delay)
		}

		CheckFrequentlyAccess(uname, 3, 2)

		// different user is not blocked.
		// uname = fmt.Sprintf("%s:%d", USER, i)

		time.Sleep(1000 * time.Millisecond)
	}
}
