package main

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"mevericcore/mcplantainer"
	"github.com/dgrijalva/jwt-go"
	"mevericcore/mcusers"
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
	jwtMdlw = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte("secret"),
		ContextKey:  "client",
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  "JWT",
		Claims:      jwt.MapClaims{},
	})
)

func main() {
	session := InitMongoDbConnection()
	defer session.Close()

	e := echo.New()

	// Debug
	e.Debug = true
	e.Logger.SetLevel(1)

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ToDo: Remove
	appG := e.Group("/app")
	appG.Use(jwtMdlw)
	mcplantainer.InitHttp(session, MainDBName, appG)
	//


	// Add Users (Auth + Me modules)
	usersG := e.Group("/users")
	mcusers.InitMainModules(session, MainDBName, usersG)

	mcplantainer.Init(session, MainDBName)

	e.Logger.Fatal(e.Start("localhost:3001"))
}
