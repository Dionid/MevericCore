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

//func (this *PlantainerDataCollectionManagerSt) SaveData(model mcmongo.ModelBaseInterface, colQuerier map[string]interface{}, data map[string]interface{}, colName string) error {
//	if err := this.UpdateModelCustomCol(colName, model, colQuerier, data); err != nil {
//		if err != this.ErrNotFound {
//			return err
//		}
//		this.InsertModelCustomCol(colName, model)
//	}
//	return nil
//}

func (this *PlantainerDataCollectionManagerSt) Init(dbsession *mgo.Session, dbName string) {
	this.AddModel(&PlantainerDataSt{})
	this.InitBase(dbsession, dbName)
	this.Inited = true
}