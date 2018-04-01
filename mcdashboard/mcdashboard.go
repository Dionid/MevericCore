package mcdashboard

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"mevericcore/mclibs/mcecho"
	"mevericcore/mclibs/mccommon"
	"mevericcore/mclibs/mcinnerrpc"
	"mevericcore/mclibs/mcws"
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8000"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true,
	}))

	e.Static("/", "")

	return e
}

var (
	innerRPCMan = mcinnerrpc.New()
)

func initInnerRPCMan() {
	innerRPCMan.Init()
	innerRPCMan.Service.Subscribe("User.RPC.Send", func(msg *mcinnerrpc.Msg) {
		rpcData := &mcws.WsRPCMsgBaseSt{}
		if err := rpcData.UnmarshalJSON(msg.Data); err != nil {
			return
		}
		WSManager.SendWsMsgByRoomName(rpcData.Dst, rpcData)
	})
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

	WSController := WSHttpControllerSt{}
	e.GET("/ws", WSController.WSHandler)

	appG := e.Group("/app")
	appG.Use(mccommon.JwtMdlw)

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

	initInnerRPCMan()

	InitUserRPCManager()

	e.Logger.Fatal(e.Start("localhost:3000"))
}