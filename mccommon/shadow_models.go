package mccommon

import "time"

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
type ShadowStateMetadataSt struct {
	Reported  *map[string]interface{}
	Desired   *map[string]interface{}
	Delta     *map[string]interface{}
	Version   int
	Timestamp time.Time
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
	Id    string `json:"id,omitempty" bson:"id"` // RequestId that will be stored on Device
	State ShadowStateSt
}

type ShadowModelInterface interface{

}