package mcdashboard

import (
	"mevericcore/mclibs/mccommon"
	"gopkg.in/mgo.v2"
)


type UsersCollectionManagerSt struct {
	mccommon.UsersCollectionManagerSt
}

func NewUsersCollectionManagerSt() *UsersCollectionManagerSt {
	return &UsersCollectionManagerSt{
		mccommon.UsersCollectionManagerSt{
			Inited: false,
		},
	}
}

func (this *UsersCollectionManagerSt) createUserAdmin() {
	email := "diodos@yandex.ru"
	admin := &UserModel{
		mccommon.UserModel{
			Login: "dionid",
			Email:    &email,
			Password: "qweqweqwe",
			IsAdmin:  true,
		},
	}
	if err := this.FindModelByLogin(admin.Login, admin); err != nil {
		if err := this.InsertModel(admin); err != nil {
			panic(err.Error())
		}
	}
}

func (this *UsersCollectionManagerSt) Init(dbsession *mgo.Session, dbName string) {
	this.AddModel(&UserModel{})
	this.InitBase(dbsession, dbName)
	this.createUserAdmin()
	this.Inited = true
}
