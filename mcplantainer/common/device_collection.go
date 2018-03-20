package common

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
)

type DevicesCollectionManagerSt struct {
	mccommon.DevicesCollectionManagerSt
	Inited bool
}

type DataCollectionManagerSt struct {
	mccommon.DataCollectionManagerSt
	Inited bool
}

func CreateNewDevicesCollectionManager(colMan mccommon.DataCollectionManagerInt) *DevicesCollectionManagerSt {
	return &DevicesCollectionManagerSt{
		mccommon.DevicesCollectionManagerSt{
			DataCollectionManager: colMan,
		},
		false,
	}
}

var (
	DataCollectionManager = &DataCollectionManagerSt{}
	DevicesCollectionManager = CreateNewDevicesCollectionManager(DataCollectionManager)
)

func InitDeviceColManager(dbsession *mgo.Session, dbName string) {
	if !DevicesCollectionManager.Inited {
		DevicesCollectionManager.AddModel(&PlantainerModelSt{})
		DevicesCollectionManager.InitManager(dbsession, dbName, "devices")
		DevicesCollectionManager.Inited = true
	}
	if !DataCollectionManager.Inited {
		DataCollectionManager.InitManager(dbsession, dbName, "plantainerdata")
		DevicesCollectionManager.Inited = true
	}
}
