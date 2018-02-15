package tztusers

import (
	"gopkg.in/mgo.v2"
	"mevericcore/mcmongo"
	"mevericcore/mccommon"
)

var UsersCollectionManager = &mccommon.UsersCollectionManager

type CompaniesCollectionManagerSt struct {
	mcmongo.CollectionManagerBaseSt
}

var CompaniesCollectionManager = CompaniesCollectionManagerSt{}

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

func initCompanyColManager(dbsession *mgo.Session, dbName string) {
	company := &mccommon.CompanyModel{}
	CompaniesCollectionManager.AddModel(company)
	CompaniesCollectionManager.InitManager(dbsession, dbName, "companies")
}

func InitCollectionsManagers(dbsession *mgo.Session, dbName string) {
	initUserColManager(dbsession, dbName)
	initCompanyColManager(dbsession, dbName)
}