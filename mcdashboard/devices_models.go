package mcdashboard

import (
	"mevericcore/mclibs/mccommon"
	"gopkg.in/mgo.v2/bson"
)

type DeviceModelSt struct {
	mccommon.DeviceWithCustomDataBaseModel `bson:",inline"`
}

//easyjson:json
type DevicesListModelSt []DeviceModelSt

func (this *DevicesListModelSt) GetBaseQuery() *bson.M {
	return &bson.M{
		"deletedAt": nil,
	}
}
