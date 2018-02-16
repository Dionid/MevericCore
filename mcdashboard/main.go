package mcdashboard

import (
	"gopkg.in/mgo.v2"
	"mevericcore/mcws"
)

var (
	WSManager = mcws.NewWSocketsManager()
)

func Init(dbsession *mgo.Session, dbName string) {
	//usersGroup, authGroup, companyGroup
	//InitUserModule(usersGroup, authGroup, companyGroup, dbsession, dbName)
}
