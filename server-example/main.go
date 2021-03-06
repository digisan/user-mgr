package main

import (
	"fmt"
	"net/http"
	"time"

	usr "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// curl test ref: https://davidwalsh.name/curl-post-file

// curl -X "POST" -F name='Foo Bar' 127.0.0.1:1323/login
// curl localhost:1323/restricted -H "Authorization: Bearer ******"

func login(c echo.Context) error {
	// [POST] Form to fill user info

	user := &usr.User{
		usr.Core{
			UName:    c.FormValue("name"),
			Email:    c.FormValue("name"),
			Password: "*pa55a@aD20TTTTT",
		},
		usr.Profile{
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
		usr.Admin{
			Regtime:   time.Now().Truncate(time.Second),
			Active:    true,
			Certified: false,
			Official:  false,
			SysRole:   "",
			MemLevel:  0,
			MemExpire: time.Time{},
			Tags:      "",
		},
	}

	fmt.Println(user)
	if user.UName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "user name is missing",
		})
	}

	claims := usr.MakeUserClaims(user)
	token := claims.GenToken()
	fmt.Println(token)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*usr.UserClaims)
	return c.String(http.StatusOK, "Welcome "+claims.UName+"!")
}

func logout(c echo.Context) error {
	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)
	claims.DeleteToken()
	return c.String(http.StatusOK, "See you "+claims.UName+"!")
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userTkn := c.Get("user").(*jwt.Token)
		claims := userTkn.Claims.(*usr.UserClaims)
		if claims.ValidateToken(userTkn.Raw) {
			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"message": "invalid or expired jwt",
		})
	}
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &usr.UserClaims{},
		SigningKey: []byte(usr.TokenKey()),
	}))
	r.Use(ValidateToken)

	r.GET("", restricted)
	r.GET("/bye", logout)

	e.Logger.Fatal(e.Start(":1323"))
}
