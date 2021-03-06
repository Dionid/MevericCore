package mccommon

import (
	"gopkg.in/mgo.v2/bson"
	"mevericcore/mclibs/mcmongo"
	"gopkg.in/mgo.v2"
)

type DevicesCollectionManagerSt struct {
	mcmongo.CollectionManagerBaseSt
	DataCollectionManager DataCollectionManagerInt
}

type DevicesCollectionManagerInterface interface {
	mcmongo.CollectionManagerBaseInterface
	SaveData(model mcmongo.ModelBaseInterface, colQuerier map[string]interface{}, data map[string]interface{}, colName string) error
	FindByOwnerId(ownerId string, modelsList DevicesWithCustomDataListBaseModelInterface) error
}

func (this *DevicesCollectionManagerSt) InitBase(dbsession *mgo.Session, dbName string) {
	this.InitManager(dbsession, dbName, "devices")
}

func (this *DevicesCollectionManagerSt) SaveData(model mcmongo.ModelBaseInterface, colQuerier map[string]interface{}, data map[string]interface{}, colName string) error {
	if err := this.UpdateModelCustomCol(colName, model, colQuerier, data); err != nil {
		if err != this.ErrNotFound {
			return err
		}
		this.InsertModelCustomCol(colName, model)
	}
	return nil
}

func (this *DevicesCollectionManagerSt) FindByOwnerId(ownerId string, modelsList DevicesWithCustomDataListBaseModelInterface) error {
	return this.FindAllModels(&bson.M{"ownersIds": bson.ObjectIdHex(ownerId)}, modelsList)
}

type DevicesWithShadowCollectionManagerSt struct {
	DevicesCollectionManagerSt
}

type DevicesWithShadowCollectionManagerInterface interface {
	DevicesCollectionManagerInterface
	FindByShadowId(shadowId string, model DeviceWithShadowBaseModelInterface) error
	DestroyByShadowId(shadowId string) error
	DeleteByShadowId(shadowId string) error
}

func (this *DevicesWithShadowCollectionManagerSt) FindByShadowId(shadowId string, model DeviceWithShadowBaseModelInterface) error {
	return this.FindModel(&bson.M{"shadow.id": shadowId}, model)
}

func (this *DevicesWithShadowCollectionManagerSt) DestroyByShadowId(shadowId string) error {
	return this.Destroy(&bson.M{"shadow.id": shadowId})
}
func (this *DevicesWithShadowCollectionManagerSt) DeleteByShadowId(shadowId string) error {
	return this.Delete(&bson.M{"shadow.id": shadowId})
}