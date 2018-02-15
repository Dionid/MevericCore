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

type ShadowStateInterface interface {

}

//easyjson:json
type ShadowModelSt struct {
	Id    string `json:"id,omitempty" bson:"id"` // RequestId that will be stored on Device
	State ShadowStateSt
}

type ShadowModelInterface interface{

}