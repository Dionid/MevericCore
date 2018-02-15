package mccommon

import (
	"gopkg.in/mgo.v2/bson"
	"mevericcore/mcmongo"
)

type DevicesCollectionManagerSt struct {
	mcmongo.CollectionManagerBaseSt
}

type DevicesCollectionManagerInterface interface {
	mcmongo.CollectionManagerBaseInterface
	SaveData(model mcmongo.ModelBaseInterface, colQuerier map[string]interface{}, data map[string]interface{}, colName string) error
	FindByOwnerId(ownerId string, modelsList DevicesListBaseModelInterface) error
	FindByShadowId(shadowId string, model DeviceBaseModelInterface) error
	DestroyByShadowId(shadowId string) error
	DeleteByShadowId(shadowId string) error
}

func (this *DevicesCollectionManagerSt) FindByOwnerId(ownerId string, modelsList DevicesListBaseModelInterface) error {
	return this.FindAllModels(&bson.M{"ownersIds": bson.ObjectIdHex(ownerId)}, modelsList)
}

func (this *DevicesCollectionManagerSt) FindByShadowId(shadowId string, model DeviceBaseModelInterface) error {
	return this.FindModel(&bson.M{"shadow.id": shadowId}, model)
}

func (this *DevicesCollectionManagerSt) DestroyByShadowId(shadowId string) error {
	return this.Destroy(&bson.M{"shadow.id": shadowId})
}
func (this *DevicesCollectionManagerSt) DeleteByShadowId(shadowId string) error {
	return this.Delete(&bson.M{"shadow.id": shadowId})
}