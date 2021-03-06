package mcplantainer

import (
	"mevericcore/mclibs/mccommon"
	"gopkg.in/mgo.v2"
	"time"
	"strconv"
	"tztatom/tztcore"
	"gopkg.in/mgo.v2/bson"
	"mevericcore/mcmodules/mclightmodule"
	"mevericcore/mcmodules/mcventilationmodule"
	"mevericcore/mcmodules/mcirrigationmodule"
)

type PlantainerCustomData struct {
	Name string `bson:"name"`
	ImageUrl string `bson:"imageUrl"`
}

type PlantainerCustomAdminData struct {

}

//easyjson:json
type PlantainerModelSt struct {
	mccommon.DeviceBaseModel `bson:",inline"`

	Shadow PlantainerShadowSt

	CustomData      PlantainerCustomData `json:"customData" bson:"customData"`
	CustomAdminData PlantainerCustomAdminData `json:"customAdminData" bson:"customAdminData"`
}

func NewPlantainerModel() *PlantainerModelSt {
	return &PlantainerModelSt{}
}

func (this *PlantainerModelSt) CheckAllSystems() (bool, error) {
	gChanged := false
	changed, err := this.Shadow.State.Reported.LightModule.CheckAllSystems(this.Shadow.State.Desired.LightModule)
	if err != nil {
		return false, err
	}
	changed, err = this.Shadow.State.Reported.IrrigationModule.CheckAllSystems(this.Shadow.State.Desired.IrrigationModule)
	if err != nil {
		return false, err
	}
	if changed {
		gChanged = true
	}
	return gChanged, nil
}

func (this *PlantainerModelSt) CheckAfterShadowReportedUpdate(oldShadow *PlantainerShadowSt) {
	this.Shadow.State.Reported.LightModule.CheckAfterShadowUpdate(this.Shadow.Id,oldShadow.State.Reported.LightModule)
	this.Shadow.State.Reported.IrrigationModule.CheckAfterShadowUpdate(this.Shadow.Id, oldShadow.State.Reported.IrrigationModule)
}

func (this *PlantainerModelSt) ExtractAndSaveData(updateData *PlantainerShadowStatePieceSt) (*PlantainerDataValuesSt, error) {
	// ToDo: Problem is that every time I save this struct, there will be empty modules structs inside
	// it's better to make method that will check if module prt of `updatedData` isn't empty and then create empty struct in `values`
	values := NewPlantainerDataValuesSt()
	addedValues := false

	if updateData.LightModule != nil {
		if values.LightModule == nil {
			values.LightModule = &mclightmodule.LightModuleStateDataSt{}
		}
		lL := updateData.LightModule.LightLvl
		if lL != nil{
			values.LightModule.LightLvl = lL
			addedValues = true
		}
		lT := updateData.LightModule.LightTurnedOn
		if lT != nil{
			values.LightModule.LightTurnedOn = lT
			addedValues = true
		}
	}

	if updateData.VentilationModule != nil {
		if values.VentilationModule == nil {
			values.VentilationModule = &mcventilationmodule.VentilationModuleStateDataSt{}
		}
		vH := updateData.VentilationModule.Humidity
		if vH != nil{
			values.VentilationModule.Humidity = vH
			addedValues = true
		}
		hCI := updateData.VentilationModule.CoolerInTurnedOn
		if hCI != nil{
			values.VentilationModule.CoolerInTurnedOn = hCI
			addedValues = true
		}
		hCO := updateData.VentilationModule.CoolerOutTurnedOn
		if hCO != nil{
			values.VentilationModule.CoolerOutTurnedOn = hCO
			addedValues = true
		}
	}

	if updateData.IrrigationModule != nil {
		if values.IrrigationModule == nil {
			values.IrrigationModule = &mcirrigationmodule.IrrigationModuleStateDataSt{}
		}

		iH := updateData.IrrigationModule.Humidity
		if iH != nil{
			values.IrrigationModule.Humidity = iH
			addedValues = true
		}

		iTO := updateData.IrrigationModule.IrrigationTurnedOn
		if iTO != nil{
			values.IrrigationModule.IrrigationTurnedOn = iTO
			addedValues = true
		}
	}

	if addedValues {
		this.CreateAndSaveData(time.Now(), values)
		return values, nil
	}

	return nil, nil
}

func (this *PlantainerModelSt) ReportedUpdate(updateData *PlantainerShadowStatePieceSt) error {
	if updateData.LightModule != nil {
		this.Shadow.State.Reported.LightModule.ReportedUpdate(updateData.LightModule)
	}
	if updateData.VentilationModule != nil {
		this.Shadow.State.Reported.VentilationModule.ReportedUpdate(&updateData.VentilationModule.VentilationModuleStateSt)
	}
	if updateData.IrrigationModule != nil {
		this.Shadow.State.Reported.IrrigationModule.ReportedUpdate(&updateData.IrrigationModule.IrrigationModuleStateSt)
	}
	// !Add new module functions here!
	return nil
}

func (this *PlantainerModelSt) DesiredUpdate(updateData *PlantainerShadowStatePieceSt) error {
	if updateData.LightModule != nil {
		this.Shadow.State.Desired.LightModule.DesiredUpdate(&updateData.LightModule.LightModuleStateSt)
	}
	if updateData.VentilationModule != nil {
		this.Shadow.State.Desired.VentilationModule.DesiredUpdate(&updateData.VentilationModule.VentilationModuleStateSt)
	}
	if updateData.IrrigationModule != nil {
		this.Shadow.State.Desired.IrrigationModule.DesiredUpdate(&updateData.IrrigationModule.IrrigationModuleStateSt)
	}
	// !Add new module functions here!
	return nil
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

func (this *PlantainerModelSt) Update(data *map[string]interface{}) error {
	println("Update")

	customDataUpdate := (*data)["customData"].(map[string]interface{})

	if customDataUpdate != nil {
		name := customDataUpdate["name"].(string)
		if name != "" {
			this.CustomData.Name = name
		}
		imageUrl := customDataUpdate["imageUrl"].(string)
		if imageUrl != "" {
			this.CustomData.ImageUrl = imageUrl
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
				&PlantainerLightModuleStateSt{},
				&PlantainerVentilationModuleStateSt{},
				&PlantainerIrrigationModuleStateSt{},
			},
			PlantainerShadowStatePieceSt{
				NewPlLightModuleStateWithDefaultsSt(),
				NewPlantainerVentilationModuleState(),
				NewPlantainerIrrigationModuleStateSt(),
			},
			nil,
		},
		PlantainerShadowMetadataSt{
			Version: 0,
		},
	}

	this.Type = this.GetTypeName()

	return nil
}

func (this *PlantainerModelSt) CreateAndSaveData(timestamp time.Time, values *PlantainerDataValuesSt) error {
	data := NewPlantainerData()

	t := time.Now()
	t.Minute()
	t.Second()
	ts := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()-t.Minute()%(data.PeriodInSec/1000), 0, 0, t.Location())
	updDataTS := t

	if !timestamp.IsZero() {
		updDataTS = timestamp
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

	if err := plantainerCollectionManager.SaveData(data, findQuery, updateQuery, "devicesData"); err != nil {
		return err
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