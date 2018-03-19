package mcusers

import (
	"mevericcore/mcws"
	"github.com/labstack/echo"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
)

var (
	WSManager = mcws.NewWSocketsManager()
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
)


func InitMainModules(dbsession *mgo.Session, dbName string, e *echo.Group) {
	adminG := e.Group("/admin")
	authG := e.Group("/auth")


	adminG.Use(jwtMdlw)

	initUsersRoutes(adminG)

	initUserModule(authG, dbsession, dbName)
	initMeModule(e)
}