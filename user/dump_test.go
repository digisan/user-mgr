package user

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/digisan/go-generics"
	fd "github.com/digisan/gotk/file-dir"
	"github.com/digisan/gotk/strs"
	lk "github.com/digisan/logkit"
	"github.com/digisan/user-mgr/db"
	. "github.com/digisan/user-mgr/user/registered"
)

const dbPath4Dump = "../data/db/user"

////
// '/registered' '/online' '/signing' '/relation' under 'dbPath4Dump' are auto generated,
//

func TestDumpList(t *testing.T) {

	db.InitDB(dbPath4Dump)
	defer db.CloseDB()

	users, err := ListUser(func(u *User) bool {
		return u.IsActive() || !u.IsActive()
	})
	lk.FailOnErr("%v", err)

	for _, user := range users {
		// registered core.go hides some fields print
		fmt.Println(user)
		fmt.Printf("%v since registered \n", user.SinceJoined())
		fmt.Println("-------------------")
	}
}

func TestDumpIngest(t *testing.T) {

	value := func(line string) string {
		v := strs.TrimHeadToFirst(line, ":")
		return strings.TrimSpace(v)
	}

	toDB := "../data/db/user"
	db.InitDB(toDB)
	defer db.CloseDB()

	u := &User{}

	fd.FileLineScan(filepath.Join(dbPath4Dump, "dump.txt"), func(line string) (bool, string) {

		switch {
		case strings.HasPrefix(line, "UName"):
			u.UName = value(line)

		case strings.HasPrefix(line, "Email"):
			u.Email = value(line)

		case strings.HasPrefix(line, "Password"):
			u.Password = value(line)

		case strings.HasPrefix(line, "Name"):
			u.Name = value(line)

		case strings.HasPrefix(line, "Active"):
			if active, ok := AnyTryToType[bool](value(line)); ok {
				u.Active = active
			} else {
				panic("Active Error")
			}

		case strings.HasPrefix(line, "RegTime"):

			v := strs.TrimTailFromFirst(value(line), ".")

			if rt, ok := AnyTryToType[time.Time](v); ok {
				u.RegTime = rt
			} else {
				println("--->", v)
				panic("RegTime Error")
			}

		case strings.HasPrefix(line, "SysRole"):
			u.SysRole = value(line)
		}

		if strs.HasAnyPrefix(line, "---") {
			fmt.Println("storing...")
			lk.FailOnErr("%v", UpdateUser(u))
			u = &User{}
		}
		return false, ""

	}, "")
}
