package main

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
)

func InitUserModule(usersGroup *echo.Group, authGroup *echo.Group, companyGroup *echo.Group, dbsession *mgo.Session, dbName string) {
	initUserColManager(dbsession, dbName)
	initUsersRoutes(usersGroup)
	initAuthRoutes(authGroup)
}
