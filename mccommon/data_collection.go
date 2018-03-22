package mccommon

import (
	"mevericcore/mcmongo"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type DataCollectionManagerSt struct {
	mcmongo.CollectionManagerBaseSt
}

type DataCollectionManagerInt interface {
	mcmongo.CollectionManagerBaseInterface
	FindByDeviceShadowId(deviceShadowId string, modelsList DeviceDataListBaseModelInterface) error
}

func (this *DataCollectionManagerSt) FindByDeviceShadowId(deviceShadowId string, modelsList DeviceDataListBaseModelInterface) error {
	return this.FindAllModels(&bson.M{"deviceShadowId": deviceShadowId}, modelsList)
}

func (this *DataCollectionManagerSt) InitBase(dbsession *mgo.Session, dbName string) {
	this.InitManager(dbsession, dbName, "devicesData")
}