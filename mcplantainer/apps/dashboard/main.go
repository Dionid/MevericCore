package main

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"mevericcore/mccommon"
	"mevericcore/mcusers"
	"github.com/labstack/echo/middleware"
	"mevericcore/mcplantainer/dashboard"
)

var (
	IsDrop = false
)

func InitMongoDbConnection() *mgo.Session {
	session, err := mgo.Dial("tzta:qweqweqwe@localhost")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	if IsDrop {
		err = session.DB("tztatom").DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	return session
}

var (
	MainDBName = "tztatom"
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

func initEcho() *echo.Echo {
	e := echo.New()

	// Debug
	e.Debug = true
	e.Logger.SetLevel(1)

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "")

	return e
}

func initUserModules(usersColMan *mccommon.UsersCollectionManagerSt, session *mgo.Session,e *echo.Echo) {
	// USERS (Auth, Me modules)
	usersG := e.Group("/users")
	mcusers.InitMain(usersColMan, usersG)
}

func initDeviceModules(e *echo.Echo) {
	devicesG := e.Group("/devices")
	dashboard.InitDeviceModule(devicesG)
}

func main() {
	// 1. Init MongoDB session
	session := InitMongoDbConnection()
	defer session.Close()

	// 2. Init Echo server for Devices and Users
	e := initEcho()

	// 3. Get UsersColManager for both modules
	usersColMan := mccommon.InitUserColManager(session, MainDBName)

	// 4. Init modules
	initUserModules(usersColMan, session, e)
	initDeviceModules(e)

	// 5. Start Echo server
	e.Logger.Fatal(e.Start("localhost:3000"))
}
