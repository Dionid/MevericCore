package mccommon

import "time"

//easyjson:json
type DeviceShadowMsgState struct {
	Reported *map[string]interface{}
	Desired  *map[string]interface{}
}

//easyjson:json
type DeviceShadowUpdateMsg struct {
	State     DeviceShadowMsgState
	Version   int
	Timestamp time.Time
}

//easyjson:json
type ShadowUpdateMsgStateSt struct {
	Reported *map[string]interface{}
	Desired  *map[string]interface{}
}

//easyjson:json
type ShadowUpdateMsgSt struct {
	State     ShadowUpdateMsgStateSt
	Version   int
	Timestamp time.Time
	ClientId  string
}

//easyjson:json
type ShadowStateMetadataTimestamp struct {
	Timestamp *time.Time
}

//easyjson:json
type ShadowStateMetadataSt struct {
	Reported  *map[string]interface{}
	Desired   *map[string]interface{}
	Delta     *map[string]interface{}
	Version   int
	Timestamp time.Time
}

func (this *ShadowStateMetadataSt) fillMetadataReported(piece *map[string]interface{}, metaReported *map[string]interface{}, now *time.Time) {
	for key, val := range *piece {
		switch v := val.(type) {
		case map[string]interface{}:
			newMap := map[string]interface{}{}
			this.fillMetadataReported(&v, &newMap, now)
			(*metaReported)[key] = newMap
		default:
			(*metaReported)[key] = ShadowStateMetadataTimestamp{
				Timestamp: now,
			}
		}
	}
}

//easyjson:json
type ShadowStateDeltaSt struct {
	Version int
	State   map[string]interface{}
}

//easyjson:json
type ShadowStateSt struct {
	Reported map[string]interface{}
	Desired  map[string]interface{}
	Delta    *ShadowStateDeltaSt
	Metadata ShadowStateMetadataSt
}

func (this *ShadowStateSt) setStatePiece(currentState *map[string]interface{}, newState *map[string]interface{}, addNew bool) {
	for key, val := range *newState {
		if (*currentState)[key] != nil {
			switch newStateV := val.(type) {
			// TODO: Test nil type
			case map[string]interface{}:
				switch currentStateV := (*currentState)[key].(type) {
				case map[string]interface{}:
					this.setStatePiece(&currentStateV, &newStateV, addNew)
				default:
					continue
				}
			default:
				(*currentState)[key] = newStateV
			}
			// TODO: This can be dangerous if system needs new values, maybe `addNew` must be always `true`
		} else if addNew {
			(*currentState)[key] = val
		}
	}
}

func (this *ShadowStateSt) SetReportedState(reported *map[string]interface{}) {
	this.setStatePiece(&this.Reported, reported, false)
}

func (this *ShadowStateSt) SetDesiredState(desired *map[string]interface{}) {
	if this.Desired == nil {
		this.Desired = map[string]interface{}{}
	}
	this.setStatePiece(&this.Desired, desired, true)
}

func (this *ShadowStateSt) CheckVersion(version int) bool {
	return this.Metadata.Version == version
}

func (this *ShadowStateSt) IncrementVersion() {
	this.Metadata.Version += 1
}

func (this *ShadowStateSt) fillDelta(reported *map[string]interface{}, desired *map[string]interface{}, delta *map[string]interface{}) {
	for key, val := range *desired {
		if (*reported)[key] != nil {
			switch desireV := val.(type) {
			case map[string]interface{}:
				switch repV := (*reported)[key].(type) {
				case map[string]interface{}:
					newMap := map[string]interface{}{}
					this.fillDelta(&repV, &desireV, &newMap)
					if len(newMap) > 0 {
						(*delta)[key] = newMap
					}
				default:
					(*delta)[key] = val
				}
			default:
				if (*reported)[key] != val {
					(*delta)[key] = val
				}
			}
		}
	}
}

func (this *ShadowStateSt) FillDelta() {
	if this.Desired == nil {
		this.Desired = map[string]interface{}{}
	}
	if this.Delta == nil {
		this.Delta = &ShadowStateDeltaSt{}
	}
	delta := map[string]interface{}{}
	this.fillDelta(&this.Reported, &this.Desired, &delta)
	this.Delta.State = delta
	this.Delta.Version = this.Metadata.Version
}

type ShadowStateInterface interface {
	FillDelta()
	GetDelta() (JSONData, bool)
}

//easyjson:json
type ShadowModelSt struct {
	Id    string `json:"id,omitempty" bson:"id"` // Id that will be stored on Device
	State ShadowStateSt
}

func (this *ShadowModelSt) GetState() *ShadowStateSt {
	return &this.State
}

type ShadowModelInterface interface{
	//ActionsOnUpdate(updateData *DeviceShadowUpdateMsg, deviceDataColMan DevicesCollectionManagerInterface) error
	GetState() *ShadowStateSt
	//NotifyOwners(msg string, handler func(userId string, msg string) (success bool))
	//SetIsActivated(colMan DevicesCollectionManagerInterface) error
}