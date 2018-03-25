package mcplantainer

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PlantainerCollectionManagerSt struct {
	mccommon.DevicesCollectionManagerSt
	Inited bool
}

func (this *PlantainerCollectionManagerSt) FindByShadowId(shadowId string, model *PlantainerModelSt) error {
	return this.FindModel(&bson.M{"shadow.id": shadowId}, model)
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