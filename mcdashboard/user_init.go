package mcdashboard

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
)

func InitUserModule(authGroup *echo.Group, dbsession *mgo.Session, dbName string) {
	initUserColManager(dbsession, dbName)
	initAuthRoutes(authGroup)
}
