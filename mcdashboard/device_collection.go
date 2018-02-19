package mcdashboard

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
)

type DevicesCollectionManagerSt struct {
	mccommon.DevicesCollectionManagerSt
}

var (
	DevicesCollectionManager = DevicesCollectionManagerSt{}
)

func initDeviceColManager(dbsession *mgo.Session, dbName string) {
	DevicesCollectionManager.AddModel(&mccommon.DeviceBaseModel{})
	DevicesCollectionManager.InitManager(dbsession, dbName, "devices")
}
