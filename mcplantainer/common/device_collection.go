package common

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
)

type PlantainerCollectionManagerSt struct {
	mccommon.DevicesCollectionManagerSt
	Inited bool
}

type PlantainerDataCollectionManagerSt struct {
	mccommon.DataCollectionManagerSt
	Inited bool
}

func CreateNewPlantainerCollectionManager(colMan mccommon.DataCollectionManagerInt) *PlantainerCollectionManagerSt {
	return &PlantainerCollectionManagerSt{
		mccommon.DevicesCollectionManagerSt{
			DataCollectionManager: colMan,
		},
		false,
	}
}

var (
	PlantainerDataCollectionManager = &PlantainerDataCollectionManagerSt{}
	PlantainerCollectionManager     = CreateNewPlantainerCollectionManager(PlantainerDataCollectionManager)
)

func InitDeviceColManager(dbsession *mgo.Session, dbName string) {
	if !PlantainerCollectionManager.Inited {
		PlantainerCollectionManager.AddModel(&PlantainerModelSt{})
		PlantainerCollectionManager.InitManager(dbsession, dbName, "devices")
		PlantainerCollectionManager.Inited = true
	}
	if !PlantainerDataCollectionManager.Inited {
		PlantainerDataCollectionManager.InitManager(dbsession, dbName, "plantainerdata")
		PlantainerCollectionManager.Inited = true
	}
}
