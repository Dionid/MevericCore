package tztusers

import (
	"gopkg.in/mgo.v2/bson"
	"tztatom/tztcoremgo"
)

import "gopkg.in/mgo.v2"

type UsersCollectionManagerSt struct {
	tztcoremgo.CollectionManagerBaseSt
}

var UsersCollectionManager = UsersCollectionManagerSt{}

func InitCollectionsManagers(dbsession *mgo.Session, dbName string) {
	admin := &UserModel{
		Email:    "diodos@yandex.ru",
		Password: "qweqweqwe",
		IsAdmin:  true,
	}
	UsersCollectionManager.AddModel(admin)
	UsersCollectionManager.InitManager(dbsession, dbName, "users")
	if err := UsersCollectionManager.FindByEmail(admin.Email, admin); err != nil {
		UsersCollectionManager.Insert(admin)
	}
}

func (this *UsersCollectionManagerSt) FindByEmail(email string, model tztcoremgo.ModelBaseInterface) error {
	return this.Find(&bson.M{"email": email}, model)
}
