package tztusers

import (
	"github.com/labstack/echo"
	"mevericcore/mcecho"
)

var (
	UserCtrl = &UserController{}
	CompanyController = &CompanyControllerSt{}
)

func InitUsersRoutes(group *echo.Group) {
	group.POST("", UserCtrl.Create)
}

func InitAuthRoutes(group *echo.Group) {
	group.POST("/login", UserCtrl.Auth)
}

func InitCompanyRoutes(group *echo.Group) {
	mcecho.CreateModelControllerRoutes(group, "/companies", CompanyController)
}
