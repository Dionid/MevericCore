package dashboard

import (
	"github.com/labstack/echo"
	"mevericcore/mcecho"
	"mevericcore/mcuserrpcmanager"
	"mevericcore/mcplantainer/common"
	"mevericcore/mccommon"
	"mevericcore/mcusers"
	"gopkg.in/mgo.v2"
)

func Init(session *mgo.Session, dbName string, e *echo.Echo) {
	common.InitDeviceColManager(session, dbName)

	plantainerColMan := common.PlantainerCollectionManager
	usersColMan := mccommon.InitUserColManager(session, dbName)

	initUserModules(usersColMan, e)
	initDeviceModules(usersColMan, e)

	appG := e.Group("/app")
	mcuserrpcmanager.InitMain(common.CreateNewPlantainerModelSt, usersColMan, plantainerColMan, appG)
}

func initUserModules(usersColMan *mccommon.UsersCollectionManagerSt, e *echo.Echo) {
	usersG := e.Group("/users")
	mcusers.InitMain(usersColMan, usersG)
}

func initDeviceModules(usersColMan *mccommon.UsersCollectionManagerSt, e *echo.Echo) {
	devicesG := e.Group("/devices")

	UserPlantainerController := &UserPlantainerControllerSt{}
	mcecho.CreateModelControllerRoutes(devicesG, "/plantainer", UserPlantainerController)
}