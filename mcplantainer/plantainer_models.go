package mcplantainer

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
	"time"
	"strconv"
	"tztatom/tztcore"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"encoding/json"
)

type PlantainerCustomData struct {
	Name string
}

type PlantainerCustomAdminData struct {

}

//easyjson:json
type PlantainerShadowStatePieceSt struct {
	LightModule PlantainerLightModuleStateSt `bson:"lightModule"`
}

func NewPlantainerShadowStatePiece() *PlantainerShadowStatePieceSt {
	return &PlantainerShadowStatePieceSt{
		*NewPlantainerLightModuleStateSt(),
	}
}

//easyjson:json
type PlantainerShadowStateSt struct {
	Reported PlantainerShadowStatePieceSt
	Desired *PlantainerShadowStatePieceSt
	Delta *PlantainerShadowStatePieceSt `bson:"-"`
}



func NewPlantainerShadowState() *PlantainerShadowStateSt {
	return &PlantainerShadowStateSt{
		*NewPlantainerShadowStatePiece(),
		nil,
		nil,
	}
}

func (this *PlantainerShadowStateSt) fillDelta(reported *map[string]interface{}, desired *map[string]interface{}, delta *map[string]interface{}) {
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
				//fmt.Printf("val : %+v\n", val)
				////isPtr := reflect.ValueOf(val).Kind() == reflect.Ptr
				//fmt.Printf("type : %+v\n", reflect.ValueOf(val).Kind())
				//fmt.Printf("type : %+v\n", reflect.ValueOf(val).Type())
				//fmt.Printf("reported : %+v\n", (*reported)[key])
				if (*reported)[key] != val && val != nil {
					(*delta)[key] = val
				}
			}
		}
	}
}

func (this *PlantainerShadowStateSt) FillDelta() *map[string]interface{} {
	des := this.Desired
	if des == nil {
		return nil
	}

	//_______

	bData, err := des.MarshalJSON()
	desMap := map[string]interface{}{}
	if err != nil {
		fmt.Printf("bData err: %+v\n", err.Error())
		return nil
	} else {
		if err := json.Unmarshal(bData, &desMap); err != nil {
			fmt.Printf("bData err: %+v\n", err.Error())
			return nil
		}
		fmt.Printf("desMap: %+v\n", desMap)
	}

	bResData, err := this.Reported.MarshalJSON()
	repMap := map[string]interface{}{}
	if err != nil {
		fmt.Printf("bResData err: %+v\n", err.Error())
		return nil
	} else {
		if err := json.Unmarshal(bResData, &repMap); err != nil {
			fmt.Printf("bResData err: %+v\n", err.Error())
			return nil
		}
		fmt.Printf("repMap: %+v\n", repMap)
	}
	deltaMap := map[string]interface{}{}
	fmt.Printf("deltaMap : %+v\n", deltaMap)
	this.fillDelta(&repMap, &desMap, &deltaMap)
	fmt.Printf("deltaMap : %+v\n", deltaMap)
	fmt.Printf("Delta: %+v\n", this.Delta)

	if len(deltaMap) == 0 {
		return nil
	}

	if dBData, err := json.Marshal(deltaMap); err != nil {
		fmt.Printf("dBData err: %+v\n", err.Error())
		return nil
	} else {
		this.Delta = &PlantainerShadowStatePieceSt{}
		fmt.Printf("dBData: %+v\n", dBData)
		if err := json.Unmarshal(dBData, &this.Delta); err != nil {
			fmt.Printf("bData err: %+v\n", err.Error())
			return nil
		}
		fmt.Printf("Success\n")
		fmt.Printf("Delta: %+v\n", this.Delta)
	}

	//desMap := structs.Map(des)
	//repMap := structs.Map(this.Reported)
	//deltaMap := map[string]interface{}{}
	//fmt.Printf("desMap : %+v\n", desMap)
	//fmt.Printf("repMap : %+v\n", repMap)
	//fmt.Printf("deltaMap : %+v\n", deltaMap)
	//this.fillDelta(&repMap, &desMap, &deltaMap)
	//fmt.Printf("deltaMap : %+v\n", deltaMap)
	//
	//if bData, err := json.Marshal(deltaMap); err != nil {
	//	fmt.Printf("bData err: %+v\n", err.Error())
	//	return nil
	//} else {
	//	if err := this.Delta.UnmarshalJSON(bData); err != nil {
	//		fmt.Printf("bData err: %+v\n", err.Error())
	//		return nil
	//	}
	//	fmt.Printf("Delta: %+v\n", this.Delta)
	//}
	//
	//if err := mapstructure.Decode(&deltaMap, &this.Delta); err != nil {
	//	fmt.Printf("err: %+v\n", err.Error())
	//	return nil
	//}

	//_______

	//rep := this.Reported

	//if this.Delta == nil {
	//	this.Delta = &PlantainerShadowStatePieceSt{}
	//}
	//
	//if rep.LightModule.LightTurnedOn != des.LightModule.LightTurnedOn {
	//	this.Delta.LightModule.LightTurnedOn = des.LightModule.LightTurnedOn
	//}
	//
	//if rep.LightModule.LightIntervalsCheckingInterval != des.LightModule.LightIntervalsCheckingInterval {
	//	this.Delta.LightModule.LightIntervalsCheckingInterval = des.LightModule.LightIntervalsCheckingInterval
	//}

	//_______

	//v := reflect.ValueOf(des)
	//values := make([]interface{}, v.NumField())
	//for i := 0; i < v.NumField(); i++ {
	//	values[i] = v.Field(i).Interface()
	//}

	//_______

	//desMap := structs.Map(des)
	//repMap := structs.Map(this.Reported)
	//delta := map[string]interface{}{}
	//this.fillDelta(&repMap, &desMap, &delta)
	//if len(delta) != 0 {
	//	if this.Delta == nil {
	//		this.Delta = &PlantainerShadowStatePieceSt{}
	//	}
	//	return &delta
	//	//if err := mapstructure.Decode(delta, this.Delta); err != nil {
	//	//	return
	//	//}
	//}
	//return nil

	//_______

	//desS := structs.New(des)
	//repS := structs.New(&this.Reported)
	//delS := structs.New(this.Delta)
	//
	//for _, f := range desS.Fields() {
	//	fmt.Printf("f name: %+v\n", f.Name())
	//	if f.IsExported() {
	//		fmt.Printf("kind : %+v\n", f.Kind())
	//		if f.Kind() == reflect.Struct {
	//			fmt.Printf("It is struct!\n")
	//			for _, ff := range f.Fields() {
	//				fmt.Printf("ff name: %+v\n", ff.Name())
	//				for _, fff := range ff.Fields() {
	//					if fff.IsExported() {
	//						if !fff.IsZero() {
	//							//fmt.Printf("field name: %+v\n", fff.Name())
	//							//fmt.Printf("value   : %+v\n", fff.Value())
	//							resF := repS.Field(f.Name()).Field(ff.Name()).Field(fff.Name())
	//							//fmt.Printf("field name: %+v\n", resF.Name())
	//							//fmt.Printf("field kind: %+v\n", resF.Kind())
	//							//fmt.Printf("value   : %+v\n", resF.Value())
	//							if resF.Kind() == reflect.Ptr {
	//								resFV := reflect.ValueOf(resF.Value()).Elem()
	//								fffFV := reflect.ValueOf(fff.Value()).Elem()
	//								fmt.Printf("%+v\n", resFV)
	//								fmt.Printf("%+v\n", fffFV)
	//								if resFV != fffFV {
	//									if _, ok := delS.Field(f.Name()); !ok {
	//
	//									}
	//								}
	//							}
	//						}
	//					}
	//				}
	//			}
	//		}
	//		//if f.Kind()
	//		//for f.Fields()
	//	}
	//}

	return nil
}

type PlantainerShadowMetadataSt struct {
	Version int
}

//easyjson:json
type PlantainerShadowSt struct {
	Id string
	State PlantainerShadowStateSt
	Metadata PlantainerShadowMetadataSt
}

func (this *PlantainerShadowSt) CheckVersion(version int) bool {
	return this.Metadata.Version == version
}

func (this *PlantainerShadowSt) IncrementVersion() {
	this.Metadata.Version += 1
}

func NewPlantainerShadow(shadowId string) *PlantainerShadowSt {
	return &PlantainerShadowSt{
		shadowId,
		*NewPlantainerShadowState(),
		PlantainerShadowMetadataSt{},
	}
}

//easyjson:json
type PlantainerModelSt struct {
	mccommon.DeviceBaseModel `bson:",inline"`

	Shadow PlantainerShadowSt

	CustomData      PlantainerCustomData `json:"customData" bson:"customData"`
	CustomAdminData PlantainerCustomAdminData `json:"customAdminData" bson:"customAdminData"`
}

func NewPlantainerModel() *PlantainerModelSt {
	return &PlantainerModelSt{
		//Shadow: PlantainerShadowSt{},
	}
}

func (this *PlantainerModelSt) EnsureIndexes(collection *mgo.Collection) error {
	return nil
}

func (this *PlantainerModelSt) GetSrc() string {
	return "/" + this.Shadow.Id
}

func (this *PlantainerModelSt) GetTypeName() string {
	return "plantainer"
}

func (this *PlantainerModelSt) CreateReported() *map[string]interface{} {
	return &map[string]interface{}{
		"lightModule": map[string]interface{}{
			"mode": "manual",
			"lightTurnedOn": false,
			"lightLvl": 0,
			"lightLvlCheckActive": false,
			"lightLvlCheckInterval": 5100,
			"lightLvlCheckLastIntervalCallTimestamp": 0,
			"lightIntervalsArr": []interface{}{},
			"lightIntervalsRestTimeTurnedOn": false,
			"lightIntervalsCheckingInterval": 20000,
		},
		"irrigationModule": map[string]interface{}{
			"mode": "manual",
			"irrigationTurnedOn": false,
			"humidity": 0,
			"temperature": 0,

			// humidityCheck
			"humidityCheckActive": false,
			"humidityCheckInterval": 5000,
			"humidityCheckLastIntervalCallTimestamp": 0,
			"humidityCheckMinLvl": 20,
			"humidityCheckAverageLvl": 25,
			"humidityCheckMaxLvl": 50,

			// irrigationTimer
			"irrigationTimerInProgress": false,
			"irrigationTimerEveryXSeconds": 11000,
			"irrigationTimerIrrigateYSeconds": 2000,
			"irrigationTimerLastCallStartTimestamp": 0,
			"irrigationTimerLastCallEndTimestamp": 0,
		},
		"airCompressorModule": map[string]interface{}{
			"turnedOn": false,
		},
		"ventilationModule": map[string]interface{}{
			"interval": 20000,
			"coolerInTurnedOn": false,
			"coolerOutTurnedOn": false,
			"humidityMaxLvl": 50,
			"humidityAverageLvl": 35,
			"mode": "manual",
		},
		"meteoStationModule": map[string]interface{}{
			"interval": 6100,
			"humidity": 0,
			"temperature": 0,
		},
	}
}

func (this *PlantainerModelSt) Update(data *map[string]interface{}) error {
	println("Update")

	customDataUpdate := (*data)["customData"].(map[string]interface{})

	if customDataUpdate != nil {
		name := customDataUpdate["name"].(string)
		if name != "" {
			this.CustomData.Name = name
		}
	}

	return nil
}

func (this *PlantainerModelSt) BeforeInsert(collection *mgo.Collection) error {

	this.CustomData.Name = "Plantainer"

	this.Shadow = PlantainerShadowSt{
		tztcore.RandString(13),
		PlantainerShadowStateSt{
			PlantainerShadowStatePieceSt{
				*NewPlantainerLightModuleStateSt(),
			},
			nil,
			nil,
		},
		PlantainerShadowMetadataSt{},
	}

	this.Type = this.GetTypeName()

	return nil
}

func (this *PlantainerModelSt) CreateAndSaveData(deviceDataColMan mccommon.DevicesCollectionManagerInterface, updateData *mccommon.DeviceShadowUpdateMsg, values *PlantainerDataValuesSt) error {
	data := NewPlantainerData()

	t := time.Now()
	t.Minute()
	t.Second()
	ts := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()-t.Minute()%(data.PeriodInSec/1000), 0, 0, t.Location())
	updDataTS := t

	if !updateData.Timestamp.IsZero() {
		updDataTS = updateData.Timestamp
	}

	minuteNum := strconv.Itoa(updDataTS.Minute() - ts.Minute())
	secondNum := strconv.Itoa(updDataTS.Second())

	//fmt.Println(updDataTS)
	//fmt.Println(ts)

	// REQUIRED FOR INSERT

	data.TS = ts
	data.DeviceShadowId = this.Shadow.Id
	data.Values = map[string]map[string]PlantainerDataValuesSt{
		minuteNum: {
			secondNum: *values,
		},
	}

	// REQUIRED FOR UPDATE
	findQuery := map[string]interface{}{"deviceShadowId": data.DeviceShadowId, "ts": data.TS}
	updateQuery := map[string]interface{}{"$set": map[string]interface{}{"values." + minuteNum + "." + secondNum: values}}

	if err := deviceDataColMan.SaveData(data, findQuery, updateQuery, "devicesData"); err != nil {
		return err
	}

	return nil
}

func (this *PlantainerModelSt) ActionsOnUpdate(updateData *mccommon.DeviceShadowUpdateMsg, deviceDataColMan mccommon.DevicesCollectionManagerInterface) error {
	println("Plantainer ActionsOnUpdate: ")

	if updateData.State.Reported != nil {
		if (*updateData.State.Reported)["lightModule"] != nil {
			//lightModuleData := (*updateData.State.Reported)["lightModule"].(map[string]interface{})
			//this.LightModule.SetState(
			//	NewLightModuleState(
			//		this.Shadow.State.Reported["mode"].(string),
			//		this.Shadow.State.Reported["lightLvlCheckActive"].(bool),
			//		int(this.Shadow.State.Reported["lightLvlCheckInterval"].(float64)),
			//		this.Shadow.State.Reported["lightIntervalsRestTimeTurnedOn"].(bool),
			//		int(this.Shadow.State.Reported["lightIntervalsCheckingInterval"].(float64)),
			//		this.Shadow.State.Reported["lightIntervalsArr"].([]LightModuleInterval),
			//	),
			//)
			//this.LightModule.CheckOnStateUpdate(
			//	this.Shadow.Id,
			//	NewLightModuleState(
			//		lightModuleData["mode"].(string),
			//		lightModuleData["lightLvlCheckActive"].(bool),
			//		int(lightModuleData["lightLvlCheckInterval"].(float64)),
			//		lightModuleData["lightIntervalsRestTimeTurnedOn"].(bool),
			//		int(lightModuleData["lightIntervalsCheckingInterval"].(float64)),
			//		lightModuleData["lightIntervalsArr"].([]LightModuleInterval),
			//	),
			//)
		}
	}

	if updateData.State.Reported != nil {
		values := NewPlantainerDataValuesSt()
		changed := false
		if (*updateData.State.Reported)["irrigationModule"] != nil {
			irrigationModuleData := (*updateData.State.Reported)["irrigationModule"].(map[string]interface{})
			if hum, ok := irrigationModuleData["humidity"].(float64); ok == true {
				hum := int(hum)
				values.IrrigationModule.Humidity = hum
				changed = true
			}
			if temperature, ok := irrigationModuleData["temperature"].(float64); ok == true {
				temperature := int(temperature)
				values.IrrigationModule.Temperature = temperature
				changed = true
			}
		}
		if (*updateData.State.Reported)["lightModule"] != nil {
			lightModuleData := (*updateData.State.Reported)["lightModule"].(map[string]interface{})
			if lightLvl, ok := lightModuleData["lightLvl"].(float64); ok == true {
				lightLvl := int(lightLvl)
				values.LightModule.LightLvl = &lightLvl
				changed = true
			}
		}
		if changed {
			this.CreateAndSaveData(deviceDataColMan, updateData, values)
		}
	}

	return nil
}

//easyjson:json
type PlantainersList []PlantainerModelSt

func (this *PlantainersList) GetBaseQuery() *bson.M {
	return &bson.M{
		"deletedAt": nil,
	}
}

func (this *PlantainersList) GetTypeName() string {
	return "plantainer"
}

func NewPlantainersList() mccommon.DevicesListBaseModelInterface {
	return &PlantainersList{}
}