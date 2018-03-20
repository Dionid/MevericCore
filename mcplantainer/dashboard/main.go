package dashboard

import (
	"github.com/labstack/echo"
	"mevericcore/mcecho"
)

func Init(e *echo.Group) {
	InitDeviceModule(e)
}

func InitDeviceModule(e *echo.Group) {
	UserPlantainerController := &UserPlantainerControllerSt{}
	mcecho.CreateModelControllerRoutes(e, "/plantainer", UserPlantainerController)
}