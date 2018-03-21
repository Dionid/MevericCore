package mccommon

import (
	"mevericcore/mcmongo"
	"time"
)

type DeviceDataBaseSt struct {
	mcmongo.ModelBase `bson:",inline"`
	TS                   time.Time `json:",omitempty" bson:"ts"`
	PeriodInSec          int       `json:"period" bson:"period"`
	DeviceShadowId       string    `json:"deviceShadowId,omitempty" bson:"deviceShadowId"`
}

type DeviceDataBaseInterface interface {
	mcmongo.ModelBaseInterface
}

//easyjson:json
type DeviceDataListBaseModelSt []DeviceDataBaseSt

type DeviceDataListBaseModelInterface interface {
	mcmongo.ModelsListBaseInterface
}