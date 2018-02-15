package main

import "github.com/labstack/echo"

var (
	UserCtrl = &UserController{}
)

func initUsersRoutes(group *echo.Group) {
	group.POST("", UserCtrl.Create)
}

func initAuthRoutes(group *echo.Group) {
	group.POST("/login", UserCtrl.Auth)
}
