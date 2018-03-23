package mcplantainer

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
)

type PlantainerCollectionManagerSt struct {
	mccommon.DevicesCollectionManagerSt
	Inited bool
}

func NewPlantainerCollectionManager(colMan mccommon.DataCollectionManagerInt) *PlantainerCollectionManagerSt {
	return &PlantainerCollectionManagerSt{
		mccommon.DevicesCollectionManagerSt{
			DataCollectionManager: colMan,
		},
		false,
	}
}

func (this *PlantainerCollectionManagerSt) Init(dbsession *mgo.Session, dbName string) {
	this.AddModel(&PlantainerModelSt{})
	this.InitBase(dbsession, dbName)
	this.Inited = true
}