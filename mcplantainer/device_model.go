package mcplantainer

import "mevericcore/mccommon"

type PlantainerCustomData struct {

}

type PlantainerCustomAdminData struct {

}

//easyjson:json
type PlantainerModelSt struct {
	mccommon.DeviceBaseModel `bson:",inline"`

	CustomData      PlantainerCustomData `json:"customData" bson:"customData"`
	CustomAdminData PlantainerCustomAdminData `json:"customAdminData" bson:"customAdminData"`
}

func CreateNewPlantainerModelSt() mccommon.DeviceBaseModelInterface {
	return &PlantainerModelSt{}
}
