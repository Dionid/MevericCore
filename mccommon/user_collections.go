package mccommon

import (
	"mevericcore/mcmongo"
	"gopkg.in/mgo.v2/bson"
)

type UsersCollectionManagerSt struct {
	mcmongo.CollectionManagerBaseSt
}

func (this *UsersCollectionManagerSt) FindModelByLogin(login string, model mcmongo.ModelBaseInterface) error {
	return this.FindModel(&bson.M{"login": login}, model)
}

func (this *UsersCollectionManagerSt) FindModelByEmail(email string, model mcmongo.ModelBaseInterface) error {
	return this.FindModel(&bson.M{"email": email}, model)
}