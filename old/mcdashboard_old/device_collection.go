package mcdashboard_old

import (
	"mevericcore/mclibs/mccommon"
	"gopkg.in/mgo.v2"
)

type DevicesCollectionManagerSt struct {
	mccommon.DevicesWithShadowCollectionManagerSt
}

var (
	DevicesCollectionManager = DevicesCollectionManagerSt{}
)

func initDeviceColManager(dbsession *mgo.Session, dbName string) {
	DevicesCollectionManager.AddModel(&mccommon.DeviceWithCustomDataBaseModel{})
	DevicesCollectionManager.InitManager(dbsession, dbName, "devices")
}
