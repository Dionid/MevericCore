package tztusers

import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
)

func InitModule(usersGroup *echo.Group, authGroup *echo.Group, dbsession *mgo.Session, dbName string) {
	InitCollectionsManagers(dbsession, dbName)
	InitUsersRoutes(usersGroup)
	InitAuthRoutes(authGroup)
}