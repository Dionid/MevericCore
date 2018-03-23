package mcplantainer

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
)

type PlantainerDataCollectionManagerSt struct {
	mccommon.DataCollectionManagerSt
	Inited bool
}

func NewPlantainerDataCollectionManager() *PlantainerDataCollectionManagerSt {
	return &PlantainerDataCollectionManagerSt{
		Inited: false,
	}
}

func (this *PlantainerDataCollectionManagerSt) Init(dbsession *mgo.Session, dbName string) {
	this.AddModel(&PlantainerDataSt{})
	this.InitBase(dbsession, dbName)
	this.Inited = true
}