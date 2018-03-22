package mcdashboard

import "mevericcore/mccommon"

type UserModel struct {
	mccommon.UserModel `bson:",inline"`
}
