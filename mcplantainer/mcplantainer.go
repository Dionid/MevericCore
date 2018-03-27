package mcplantainer


import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	isDBDrop = false
	mainDBName = "tztatom"
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
	plantainerDataCollectionManager = NewPlantainerDataCollectionManager()
	plantainerCollectionManager = NewPlantainerCollectionManager(plantainerDataCollectionManager)
)

func initCollections(session *mgo.Session) {
	plantainerDataCollectionManager.Init(session, mainDBName)
	plantainerCollectionManager.Init(session, mainDBName)
}

var (
	deviceCronManager = NewDeviceCronManager()
)

func Init() {
	// 1. Init MongoDB session
	session := initMongoDbConnection()
	defer session.Close()

	// 2. Init Echo server for Devices and Users
	e := initEcho()

	// 3. Init Collections
	initCollections(session)
	initUserRPCManager()
	initDeviceRPCManager()
	initCronRPCManager()
	initInnerRPCMan()

	deviceCronManager.Init()

	initMQTT()

	e.Logger.Fatal(e.Start("localhost:3001"))
}
