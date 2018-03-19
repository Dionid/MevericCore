package mcusers

import (
	"mevericcore/mcmongo"
	"gopkg.in/mgo.v2/bson"
)

//easyjson:json
type MeModel struct {
	mcmongo.ModelBase `bson:",inline"`
	Login  string  `bson:",omitempty"`
	Email                *string  `bson:",omitempty"`
	Phone                *string `json:",omitempty"`
	CompanyId  *bson.ObjectId `json:",omitempty" bson:"companyId"`
}

