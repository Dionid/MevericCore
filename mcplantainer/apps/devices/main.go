package main

import (
	"gopkg.in/mgo.v2"
	"mevericcore/mcplantainer/device"
)

var (
	IsDrop = false
)

func InitMongoDbConnection() *mgo.Session {
	session, err := mgo.Dial("tzta:qweqweqwe@localhost")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	if IsDrop {
		err = session.DB("tztatom").DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	return session
}

var (
	MainDBName = "tztatom"
)

func main() {
	// 1. Init MongoDB session
	session := InitMongoDbConnection()
	defer session.Close()

	// 4. Init modules
	device.Init(session, MainDBName)

	select {}
}
