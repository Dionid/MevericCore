package mccommon

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"mevericcore/mcmongo"
	"errors"
)

//easyjson:json
type DeviceBaseModel struct {
	mcmongo.ModelBase `bson:",inline"`
	Shadow               ShadowModelSt `bson:"shadow"`

	//Name string `json:"name,omitempty" bson:"name"`
	Src  string `json:"srcId,omitempty" bson:"src"` // Channel Id where this device is working
	Type string `json:"type,omitempty" bson:"type"`

	FirstActivation *time.Time `json:"firstActivation,omitempty" bson:"firstActivation,omitempty"` // Activation date
	LastSeenOnline  *time.Time `json:"lastSeenOnline,omitempty" bson:"lastSeenOnline,omitempty"`
	IsOnline        *bool      `json:"isOnline" bson:"isOnline"` // Is now online and working

	OwnersIds []bson.ObjectId          `json:"ownersIds" bson:"ownersIds,omitempty"`
	Owners    *UsersListModel `json:"-" bson:"-"`

	// ToDo: Add this
	// NetworkId []string // Representation of network for this device

	CustomData      map[string]interface{} `json:"customData" bson:"customData"`
	CustomAdminData map[string]interface{} `json:"customAdminData" bson:"customAdminData"`
}

type DeviceBaseModelInterface interface {
	mcmongo.ModelBaseInterface
	MarshalJSON() ([]byte, error)
	UpdateCustomData(data *map[string]interface{}) bool
	UpdateCustomAdminData(data *map[string]interface{}) bool
	GetTypeName() string
}

func (this *DeviceBaseModel) GetTypeName() string {
	return ""
}

func (this *DeviceBaseModel) ActionsOnUpdate(updateData *ShadowUpdateMsgSt, deviceDataColMan DevicesCollectionManagerInterface) error {
	return nil
}

func (this *DeviceBaseModel) GetState() ShadowStateInterface {
	return nil
}

func (this *DeviceBaseModel) GetSrc() string {
	return ""
}

func (this *DeviceBaseModel) NotifyOwners(msg string, handler func(userId string, msg string) (success bool)) {
	return
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

func (this *DeviceBaseModel) updateCustomData(currentState *map[string]interface{}, newState *map[string]interface{}) bool {
	for key, val := range *newState {
		if (*currentState)[key] != nil {
			switch newStateV := val.(type) {
			// TODO: Test nil type
			case nil:
				delete(*currentState, key)
				continue
			case map[string]interface{}:
				switch currentStateV := (*currentState)[key].(type) {
				case map[string]interface{}:
					this.updateCustomData(&currentStateV, &newStateV)
				default:
					(*currentState)[key] = newStateV
					continue
				}
			default:
				(*currentState)[key] = newStateV
			}
			// TODO: This can be dangerous if system needs new values, maybe `addNew` must be always `true`
		} else {
			switch val.(type) {
			case nil:
				continue
			default:
				(*currentState)[key] = val
			}
		}
	}
	return true
}

func (this *DeviceBaseModel) UpdateCustomAdminData(data *map[string]interface{}) bool {
	return this.updateCustomData(&this.CustomAdminData, data)
}

func (this *DeviceBaseModel) UpdateCustomData(data *map[string]interface{}) bool {
	return this.updateCustomData(&this.CustomData, data)
}

//easyjson:json
type DevicesListBaseModel []DeviceBaseModel

type DevicesListBaseModelInterface interface {
	mcmongo.ModelsListBaseInterface
	GetTypeName() string
}

func (this *DevicesListBaseModel) GetBaseQuery() *bson.M {
	return &bson.M{
		"deletedAt": nil,
	}
}

func (this *DevicesListBaseModel) GetTypeName() string {
	return ""
}