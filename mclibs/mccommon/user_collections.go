package mccommon

import (
	"mevericcore/mclibs/mcmongo"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type UsersCollectionManagerSt struct {
	mcmongo.CollectionManagerBaseSt
	Inited bool
}

func (this *UsersCollectionManagerSt) FindModelByLogin(login string, model mcmongo.ModelBaseInterface) error {
	return this.FindModel(&bson.M{"login": login}, model)
}

func (this *UsersCollectionManagerSt) FindModelByEmail(email string, model mcmongo.ModelBaseInterface) error {
	return this.FindModel(&bson.M{"email": email}, model)
}

func (this *UsersCollectionManagerSt) InitBase(dbsession *mgo.Session, dbName string) {
	this.InitManager(dbsession, dbName, "users")
}

func createUserAdmin() {
	email := "diodos@yandex.ru"
	admin := &UserModel{
		Login: "dionid",
		Email:    &email,
		Password: "qweqweqwe",
		IsAdmin:  true,
	}
	if err := UsersCollectionManager.FindModelByLogin(admin.Login, admin); err != nil {
		if err := UsersCollectionManager.InsertModel(admin); err != nil {
			panic(err.Error())
		}
	}
}

var UsersCollectionManager = &UsersCollectionManagerSt{}

func InitUserColManager(dbsession *mgo.Session, dbName string) *UsersCollectionManagerSt {
	if !UsersCollectionManager.Inited {
		UsersCollectionManager.AddModel(&UserModel{})
		UsersCollectionManager.InitManager(dbsession, dbName, "users")
		createUserAdmin()
		UsersCollectionManager.Inited = true
	}
	return UsersCollectionManager
}