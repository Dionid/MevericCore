package mcdashboard

import "github.com/labstack/echo"

var (
	MeCtrl = &MeController{}
)

func initMeRoutes(group *echo.Group) {
	group.GET("", MeCtrl.Me)
}
