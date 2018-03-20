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
	DataCollectionManager       = &DataCollectionManagerSt{}
	PlantainerCollectionManager = CreateNewDevicesCollectionManager(DataCollectionManager)
)

func InitDeviceColManager(dbsession *mgo.Session, dbName string) {
	if !PlantainerCollectionManager.Inited {
		PlantainerCollectionManager.AddModel(&PlantainerModelSt{})
		PlantainerCollectionManager.InitManager(dbsession, dbName, "devices")
		PlantainerCollectionManager.Inited = true
	}
	if !DataCollectionManager.Inited {
		DataCollectionManager.InitManager(dbsession, dbName, "plantainerdata")
		PlantainerCollectionManager.Inited = true
	}
}
