package main

import (
	"gopkg.in/mgo.v2"
	"mevericcore/mcplantainer/device"
)

var (
	IsDrop = false
	MainDBName = "tztatom"
)

func InitMongoDbConnection() *mgo.Session {
	session, err := mgo.Dial("tzta:qweqweqwe@localhost")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	if IsDrop {
		err = session.DB(MainDBName).DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	return session
}

func main() {
	// 1. Init MongoDB session
	session := InitMongoDbConnection()
	defer session.Close()

	// 4. Init modules
	device.Init(session, MainDBName)

	println("App activated")

	select {}
}
