package mcusers

import (
	"mevericcore/mcws"
	"github.com/labstack/echo"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
	"mevericcore/mccommon"
)

var (
	WSSkipperFn = func(c echo.Context) bool {
		fmt.Println(c.Path())
		if c.Path() == "/app/ws" {
			return true
		}
		return false
	}
	jwtMdlw = middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper:     WSSkipperFn,
		SigningKey:  []byte("secret"),
		ContextKey:  "client",
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  "JWT",
		Claims:      jwt.MapClaims{},
	})
	UsersCollectionManager *mccommon.UsersCollectionManagerSt = nil
)

func InitColManagers(userColManager *mccommon.UsersCollectionManagerSt) {
	UsersCollectionManager = userColManager
}

func InitMainHttp(e *echo.Group) {
	adminG := e.Group("/admin")
	adminG.Use(jwtMdlw)
	initUsersRoutes(adminG)

	authG := e.Group("/auth")
	initAuthRoutes(authG)
	initMeModule(e)
}

func InitMain(userColManager *mccommon.UsersCollectionManagerSt, e *echo.Group) {
	InitColManagers(userColManager)
	InitMainHttp(e)
}