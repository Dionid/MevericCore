package mcdashboard

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
)

type DevicesCollectionManagerSt struct {
	mccommon.DevicesWithShadowCollectionManagerSt
	Inited bool
}

func NewDeviceCollectionManager(colMan mccommon.DataCollectionManagerInt) *DevicesCollectionManagerSt {
	return &DevicesCollectionManagerSt{
		mccommon.DevicesWithShadowCollectionManagerSt{
			DevicesCollectionManagerSt: mccommon.DevicesCollectionManagerSt{
				DataCollectionManager: colMan,
			},
		},
		false,
	}
}

func (this *DevicesCollectionManagerSt) Init(dbsession *mgo.Session, dbName string) {
	this.AddModel(&DeviceModelSt{})
	this.InitBase(dbsession, dbName)
	this.Inited = true
}

type DeviceDataCollectionManagerSt struct {
	mccommon.DataCollectionManagerSt
	Inited bool
}

func (this *DeviceDataCollectionManagerSt) Init(dbsession *mgo.Session, dbName string) {
	this.AddModel(&DeviceModelSt{})
	this.InitBase(dbsession, dbName)
	this.Inited = true
}

func NewDeviceDataCollectionManager() *DeviceDataCollectionManagerSt {
	return &DeviceDataCollectionManagerSt{

	}
}