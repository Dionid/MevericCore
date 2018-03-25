package mcdashboard_old

import (
	"gopkg.in/mgo.v2"
)

func InitDeviceModule(dbsession *mgo.Session, dbName string) {
	initDeviceColManager(dbsession, dbName)
}