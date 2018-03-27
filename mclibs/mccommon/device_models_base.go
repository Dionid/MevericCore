package mccommon

import (
	"mevericcore/mclibs/mcmongo"
	"time"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"gopkg.in/mgo.v2"
)

//easyjson:json
type DeviceBaseModel struct {
	mcmongo.ModelBase `bson:",inline"`

	Src  string `json:"srcId,omitempty" bson:"src"` // Channel Id where this device is working
	Type string `json:"type,omitempty" bson:"type"`

	FirstActivation *time.Time `json:"firstActivation,omitempty" bson:"firstActivation,omitempty"` // Activation date
	LastSeenOnline  *time.Time `json:"lastSeenOnline,omitempty" bson:"lastSeenOnline,omitempty"`
	IsOnline        *bool      `json:"isOnline" bson:"isOnline"` // Is now online and working

	OwnersIds []bson.ObjectId          `json:"ownersIds" bson:"ownersIds,omitempty"`
	Owners    *UsersListModel `json:"-" bson:"-"`
}

type DeviceBaseModelInterface interface {
	mcmongo.ModelBaseInterface
	IsOwnerStringId(ownerId string) (bool, error)
	IsOwner(ownerId bson.ObjectId) (bool, error)
	MarshalJSON() ([]byte, error)
	GetTypeName() string
	Update(*map[string]interface{}) error
}

func (this *DeviceBaseModel) EnsureIndex(collection *mgo.Collection) error {
	return nil
}

func (this *DeviceBaseModel) Update(data *map[string]interface{}) error {
	return nil
}

func (this *DeviceBaseModel) GetTypeName() string {
	return ""
}

func (this *DeviceBaseModel) IsOwner(ownerId bson.ObjectId) (bool, error) {

	if this.OwnersIds == nil {
		return false, errors.New("owners must be fulfilled")
	}

	for _, id := range this.OwnersIds {
		if ownerId == id {
			return true, nil
		}
	}

	return false, nil
}

func (this *DeviceBaseModel) IsOwnerStringId(ownerId string) (bool, error) {
	return this.IsOwner(bson.ObjectIdHex(ownerId))
}
