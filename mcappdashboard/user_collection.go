package main

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
)

type UsersCollectionManagerSt struct {
	mccommon.UsersCollectionManagerSt
}

var (
	UsersCollectionManager = UsersCollectionManagerSt{}
)

func createUserAdmin() {
	email := "diodos@yandex.ru"
	admin := &mccommon.UserModel{
		Login: "dionid",
		Email:    &email,
		Password: "qweqweqwe",
		IsAdmin:  true,
	}
	if err := UsersCollectionManager.FindModelByLogin(admin.Login, admin); err != nil {
		UsersCollectionManager.InsertModel(admin)
	}
}

func initUserColManager(dbsession *mgo.Session, dbName string) {
	UsersCollectionManager.AddModel(&mccommon.UserModel{})
	UsersCollectionManager.InitManager(dbsession, dbName, "users")
	createUserAdmin()
}