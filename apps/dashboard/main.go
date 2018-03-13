package main

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"mevericcore/mcdashboard"
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
)

func main() {
	session := InitMongoDbConnection()
	defer session.Close()

	e := echo.New()

	// Debug
	e.Debug = true
	e.Logger.SetLevel(1)

	// Middleware
	//e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//e.Static("/public", "./public")
	//e.File("./public/index.html", "")
	//e.File("", "public/index.html")

	e.GET("/alina", func(c echo.Context) error {
		return c.File("./public/index.html")
	})

	mcdashboard.Init(session, MainDBName, e)

	e.Logger.Fatal(e.Start("localhost:3000"))
}
