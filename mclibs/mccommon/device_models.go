package mccommon

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"mevericcore/mclibs/mcmongo"
)

//easyjson:json
type DeviceWithShadowBaseModel struct {
	DeviceBaseModel `bson:",inline"`
	Shadow               ShadowModelSt `bson:"shadow"`
}

type DeviceWithShadowBaseModelInterface interface {
	DeviceBaseModelInterface
	GetShadow() ShadowModelInterface
	GetShadowId() string
	ActionsOnUpdate(updateData *DeviceShadowUpdateMsg, deviceDataColMan DevicesWithShadowCollectionManagerInterface) error
}

type DeviceCreatorFn func() DeviceWithShadowBaseModelInterface
type DevicesListCreatorFn func() DevicesListBaseModelInterface

func (this *DeviceWithShadowBaseModel) GetShadow() ShadowModelInterface {
	return &this.Shadow
}

func (this *DeviceWithShadowBaseModel) GetShadowId() string {
	return this.Shadow.Id
}

func (this *DeviceWithShadowBaseModel) CreateShadowStateMetadata(reported *map[string]interface{}) *ShadowStateMetadataSt {
	now := time.Now()

	state := &ShadowStateMetadataSt{
		Version:   0,
		Timestamp: now,
		Reported:  &map[string]interface{}{},
	}

	state.fillMetadataReported(reported, state.Reported, &now)

	return state
}

func (this *DeviceWithShadowBaseModel) CreateShadowState(reported *map[string]interface{}) *ShadowStateSt {
	return &ShadowStateSt{
		Reported: *reported,
		Desired:  nil,
		Delta:    nil,
		Metadata: *this.CreateShadowStateMetadata(reported),
	}
}

func (this *DeviceWithShadowBaseModel) ActionsOnUpdate(updateData *DeviceShadowUpdateMsg, deviceDataColMan DevicesWithShadowCollectionManagerInterface) error {
	return nil
}

func (this *DeviceWithShadowBaseModel) GetState() ShadowStateInterface {
	return nil
}

func (this *DeviceWithShadowBaseModel) GetSrc() string {
	return ""
}

func (this *DeviceWithShadowBaseModel) NotifyOwners(msg string, handler func(userId string, msg string) (success bool)) {
	return
}

//easyjson:json
type DeviceWithCustomDataBaseModel struct {
	DeviceWithShadowBaseModel `bson:",inline"`

	CustomData      map[string]interface{} `json:"customData" bson:"customData"`
	CustomAdminData map[string]interface{} `json:"customAdminData" bson:"customAdminData"`
}

type DeviceWithCustomDataBaseModelInterface interface {
	DeviceWithShadowBaseModelInterface
	UpdateCustomData(data *map[string]interface{}) bool
	UpdateCustomAdminData(data *map[string]interface{}) bool
}

func (this *DeviceWithCustomDataBaseModel) updateCustomData(currentState *map[string]interface{}, newState *map[string]interface{}) bool {
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

func (this *DeviceWithCustomDataBaseModel) UpdateCustomAdminData(data *map[string]interface{}) bool {
	return this.updateCustomData(&this.CustomAdminData, data)
}

func (this *DeviceWithCustomDataBaseModel) UpdateCustomData(data *map[string]interface{}) bool {
	return this.updateCustomData(&this.CustomData, data)
}

//easyjson:json
type DevicesListBaseModel []DeviceWithShadowBaseModel

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


//easyjson:json
type DevicesWithCustomDataListBaseModel []DeviceWithCustomDataBaseModel

type DevicesWithCustomDataListBaseModelInterface interface {
	mcmongo.ModelsListBaseInterface
}

func (this *DevicesWithCustomDataListBaseModel) GetBaseQuery() *bson.M {
	return &bson.M{
		"deletedAt": nil,
	}
}

func (this *DevicesWithCustomDataListBaseModel) GetTypeName() string {
	return ""
}