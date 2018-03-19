package mcusers

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
)

func initUserModule(authGroup *echo.Group, dbsession *mgo.Session, dbName string) {
	initUserColManager(dbsession, dbName)
	initAuthRoutes(authGroup)
}
