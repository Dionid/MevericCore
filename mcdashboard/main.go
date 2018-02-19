package mcdashboard

import (
	"gopkg.in/mgo.v2"
	"mevericcore/mcws"
	"github.com/labstack/echo"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
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
		ContextKey:  "user",
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  "JWT",
		Claims:      jwt.MapClaims{},
	})
)


func Init(dbsession *mgo.Session, dbName string, e *echo.Echo) {
	adminG := e.Group("/admin")
	adminG.Use(jwtMdlw)
	adminUserG := adminG.Group("/users")

	authG := e.Group("/auth")

	appG := e.Group("/app")
	appG.Use(jwtMdlw)

	//deviceG := mcecho.CreateModelControllerRoutes(appG, "/device", )

	InitUserModule(adminUserG, authG, dbsession, dbName)
	initWsRoute(appG)
	InitWsManager()
}
