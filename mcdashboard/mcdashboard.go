package mcdashboard

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
	"mevericcore/mcecho"
)

var (
	isDBDrop = false
	mainDBName = "tztatom"
	wsSkipperFn = func(c echo.Context) bool {
		fmt.Println(c.Path())
		if c.Path() == "/app/ws" {
			return true
		}
		return false
	}
	jwtMdlw = middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper:     wsSkipperFn,
		SigningKey:  []byte("secret"),
		ContextKey:  "client",
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  "JWT",
		Claims:      jwt.MapClaims{},
	})
)

func initMongoDbConnection() *mgo.Session {
	session, err := mgo.Dial("tzta:qweqweqwe@localhost")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	if isDBDrop {
		err = session.DB("tztatom").DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	return session
}

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

var (
	usersCollectionManager = NewUsersCollectionManagerSt()
	devicesDataCollectionManager = NewDeviceDataCollectionManager()
	devicesCollectionManager = NewDeviceCollectionManager(devicesDataCollectionManager)
)

func initCollections(session *mgo.Session) {
	usersCollectionManager.Init(session, mainDBName)
	devicesDataCollectionManager.Init(session, mainDBName)
	devicesCollectionManager.Init(session, mainDBName)
}

func initRoutes(e *echo.Echo) {
	authG := e.Group("/auth")
	initAuthRoutes(authG)

	appG := e.Group("/app")
	appG.Use(jwtMdlw)

	meG := appG.Group("/me")
	initMeRoutes(meG)

	mcecho.CreateModelControllerRoutes(appG, "/devices", &UserDevicesControllerSt{})
}

func Init() {
	// 1. Init MongoDB session
	session := initMongoDbConnection()
	defer session.Close()

	// 2. Init Echo server for Devices and Users
	e := initEcho()

	// 3. Init Collections
	initCollections(session)

	// 4. Init Routes
	initRoutes(e)

	e.Logger.Fatal(e.Start("localhost:3000"))
}