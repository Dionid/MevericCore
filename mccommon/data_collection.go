package mccommon

import (
	"mevericcore/mcmongo"
	"gopkg.in/mgo.v2/bson"
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