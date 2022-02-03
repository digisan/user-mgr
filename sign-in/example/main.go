package main

import (
	"fmt"

	lk "github.com/digisan/logkit"
	si "github.com/digisan/user-mgr/sign-in"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
)

func main() {

	udb.OpenSession("../../data/user")
	defer udb.CloseSession()

	// get [user] from GET

	// Will be GET header
	user := usr.User{
		UName:    "QMiao",
		Password: "pa55a@aD20TTTTT",
	}

	lk.FailOnErr("%v", si.UserExists(user))
	lk.FailOnErrWhen(!si.PwdOK(user), "%v", fmt.Errorf("incorrect password"))
	lk.FailOnErr("%v", si.Trail(user.UName))

	fmt.Println("Login OK")
}
