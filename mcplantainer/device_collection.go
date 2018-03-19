package mcplantainer

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
)

type DevicesCollectionManagerSt struct {
	mccommon.DevicesCollectionManagerSt
}

type DataCollectionManagerSt struct {
	mccommon.DataCollectionManagerSt
}

func CreateNewDevicesCollectionManager(colMan mccommon.DataCollectionManagerInt) *DevicesCollectionManagerSt {
	return &DevicesCollectionManagerSt{
		mccommon.DevicesCollectionManagerSt{
			DataCollectionManager: colMan,
		},
	}
}

var (
	DataCollectionManager = &DataCollectionManagerSt{}
	DevicesCollectionManager = CreateNewDevicesCollectionManager(DataCollectionManager)
)

func initDeviceColManager(dbsession *mgo.Session, dbName string) {
	DevicesCollectionManager.AddModel(&PlantainerModelSt{})
	DevicesCollectionManager.InitManager(dbsession, dbName, "devices")
	DataCollectionManager.InitManager(dbsession, dbName, "plantainerdata")
}
