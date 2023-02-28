package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	lk "github.com/digisan/logkit"
	si "github.com/digisan/user-mgr/sign-in"
	so "github.com/digisan/user-mgr/sign-out"
	usr "github.com/digisan/user-mgr/user"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// curl test ref: https://davidwalsh.name/curl-post-file

// curl -X "POST" -F name='Foo Bar' 127.0.0.1:1323/login
// curl localhost:1323/restricted -H "Authorization: Bearer ******"

var user = &usr.User{
	Core: usr.Core{
		UName:    "",
		Email:    "",
		Password: "*pa55a@aD20TTTTT",
	},
	Profile: usr.Profile{
		Name:           "",
		Phone:          "",
		Country:        "",
		City:           "",
		Addr:           "",
		PersonalIDType: "",
		PersonalID:     "",
		Gender:         "",
		DOB:            "",
		Position:       "",
		Title:          "",
		Employer:       "",
		Bio:            "",
		AvatarType:     "",
		Avatar:         []byte{},
	},
	Admin: usr.Admin{
		RegTime:   time.Now().Truncate(time.Second),
		Active:    false,
		Certified: false,
		Official:  false,
		SysRole:   "",
		MemLevel:  0,
		MemExpire: time.Time{},
		Tags:      "",
	},
}

func login(c echo.Context) error {
	// [POST] Form to fill user info

	lk.Log("login")

	user.UName = c.FormValue("name")
	user.Email = c.FormValue("name")

	fmt.Println(user)
	if user.UName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "user name is missing",
		})
	}

	claims := usr.MakeClaims(user)
	token := usr.GenerateToken(claims)
	fmt.Println(token)

	lk.FailOnErr("%v", si.UserStatusIssue(user))                              // check user existing status
	lk.FailOnErrWhen(!si.PwdOK(user), "%v", fmt.Errorf("incorrect password")) // check password
	lk.FailOnErr("%v", si.Trail(user.UName))                                  // this is a user online record notification

	fmt.Println("Login OK")

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func accessible(c echo.Context) error {

	lk.Log("accessible")

	return c.String(http.StatusOK, "Accessible")
}

func auth(c echo.Context) error {

	lk.Log("auth")

	invoker, err := usr.Invoker(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Welcome "+invoker.UName+"!")
}

func logout(c echo.Context) error {

	lk.Log("logout")

	invoker, err := usr.Invoker(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	invoker.DeleteToken()
	return c.String(http.StatusOK, "See you "+invoker.UName+"!")
}

func activate(c echo.Context) error {

	lk.Log("activate")

	invoker, err := usr.Invoker(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	invoker.Active = true
	return c.String(http.StatusOK, invoker.UName+" is activated!")
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, claims, err := usr.TokenClaimsInHandler(c)
		if err != nil {
			return err
		}
		invoker := usr.ClaimsToUser(claims)
		if invoker.ValidateToken(token.Raw) {
			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"message": "invalid or expired jwt",
		})
	}
}

func main() {

	usr.InitDB("./data/user")
	defer usr.CloseDB()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	///////////////////////////////////////////////////////

	offline := make(chan string, 2048)
	si.SetOfflineTimeout(10 * time.Second)
	si.MonitorOffline(ctx, offline, func(uname string) error { return so.Logout(uname) })
	go func() {
		for rm := range offline {
			fmt.Println("offline:", rm)
			if e := so.Logout(rm); e != nil {
				log.Fatalf("offline error @%s on %v", rm, e)
			}
		}
	}()

	///////////////////////////////////////////////////////

	usr.SetTokenValidPeriod(20 * time.Second)
	usr.MonitorTokenExpired(ctx, func(uname string) error {
		fmt.Printf("[%s]'s session is expired\n", uname)
		return nil
	})

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/", accessible)

	// Auth group
	r := e.Group("/auth")

	// Configure middleware with the custom claims type
	r.Use(echojwt.JWT([]byte(usr.TokenKey())))
	r.Use(ValidateToken)

	r.GET("", auth)
	r.GET("/bye", logout)
	r.POST("/activate", activate)

	e.Logger.Fatal(e.Start(":1323"))
}
