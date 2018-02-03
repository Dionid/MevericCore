package tztusers

import "github.com/labstack/echo"

var UserCtrl = &UserController{}

func InitUsersRoutes(group *echo.Group) {
	group.POST("", UserCtrl.Create)
}

func InitAuthRoutes(group *echo.Group) {
	group.POST("/login", UserCtrl.Auth)
}
