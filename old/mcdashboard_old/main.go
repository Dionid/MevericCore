package mcdashboard_old

import (
	"gopkg.in/mgo.v2"
	"mevericcore/mcws"
	"github.com/labstack/echo"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
	"mevericcore/mcecho"
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


func Init(dbsession *mgo.Session, dbName string, e *echo.Echo) {
	adminG := e.Group("/admin")
	appG := e.Group("/app")
	authG := e.Group("/auth")


	adminG.Use(jwtMdlw)
	appG.Use(jwtMdlw)

	initUsersRoutes(adminG)

	mcecho.CreateModelControllerRoutes(appG, "/devices", DeviceCtrl)

	InitUserModule(authG, dbsession, dbName)
	InitDeviceModule(dbsession, dbName)

	// 1.1. Me
	initMeModule(e)

	initWsRoute(appG)
	InitWsManager()
}
