package common

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"tztatom/tztcore"
	"time"
	"strconv"
	"fmt"
	"mevericcore/mcmongo"
)

type PlantainerDataValuesIrrigationModuleSt struct {
	Humidity int
}

type PlantainerDataValuesSt struct {
	IrrigationModule *PlantainerDataValuesIrrigationModuleSt
}

func NewPlantainerDataValuesSt() *PlantainerDataValuesSt {
	return &PlantainerDataValuesSt{
		&PlantainerDataValuesIrrigationModuleSt{},
	}
}

type PlantainerDataSt struct {
	mcmongo.ModelBase `bson:",inline"`
	TS                   time.Time `json:",omitempty" bson:"ts"`
	PeriodInSec          int       `json:"period" bson:"period"`
	DeviceShadowId       string    `json:"deviceShadowId,omitempty" bson:"deviceShadowId"`
	Values               map[string]map[string]PlantainerDataValuesSt
}

func NewPlantainerData() *PlantainerDataSt {
	return &PlantainerDataSt{
		PeriodInSec: 10000,
	}
}

type PlantainerCustomData struct {
	Name string
}

type PlantainerCustomAdminData struct {

}

//easyjson:json
type PlantainerModelSt struct {
	mccommon.DeviceBaseModel `bson:",inline"`

	CustomData      PlantainerCustomData `json:"customData" bson:"customData"`
	CustomAdminData PlantainerCustomAdminData `json:"customAdminData" bson:"customAdminData"`
}

func CreateNewPlantainerModelSt() mccommon.DeviceBaseModelInterface {
	return &PlantainerModelSt{}
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
		},
	}
}

func (this *PlantainerModelSt) BeforeInsert(collection *mgo.Collection) error {

	this.CustomData.Name = "Plantainer"

	if this.Shadow.Id == "" {
		// TODO: test shadowId
		this.Shadow.Id = tztcore.RandString(13)
	}

	reported := this.CreateReported()

	this.Shadow.State = *this.CreateShadowState(reported)
	this.Type = this.GetTypeName()

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
		values := NewPlantainerDataValuesSt()
		if (*updateData.State.Reported)["irrigationModule"] != nil {
			irrigationModuleData := (*updateData.State.Reported)["irrigationModule"].(map[string]interface{})
			if hum, ok := irrigationModuleData["humidity"].(float64); ok == true {
				hum := int(hum)
				values.IrrigationModule.Humidity = hum
			}
		}
		if *values != (PlantainerDataValuesSt{}) {
			this.CreateAndSaveData(deviceDataColMan, updateData, values)
		}
	}

	return nil
}