package mcplantainer

import "mevericcore/mccommon"

type PlantainerCustomData struct {

}

type PlantainerCustomAdminData struct {

}

//easyjson:json
type PlantainerModelSt struct {
	mccommon.DeviceBaseModel

	CustomData      PlantainerCustomData `json:"customData" bson:"customData"`
	CustomAdminData PlantainerCustomAdminData `json:"customAdminData" bson:"customAdminData"`
}
