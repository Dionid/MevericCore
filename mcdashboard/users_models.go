package mcdashboard

import "mevericcore/mclibs/mccommon"

type UserModel struct {
	mccommon.UserModel `bson:",inline"`
}
