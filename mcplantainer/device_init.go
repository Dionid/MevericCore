package mcplantainer

import (
	"mevericcore/mcdevicerpcmanager"
)

type PlantainerCtrlSt struct {
	mcdevicerpcmanager.DeviceRPCCtrlSt
}

func CreateNewPlantainerCtrl(typeName string) *PlantainerCtrlSt {
	bR := mcdevicerpcmanager.CreateNewDeviceRPCCtrl(typeName, CreateNewPlantainerModelSt, DeviceMQTTManager)
	bR.DevicesCollectionManager = DevicesCollectionManager

	res := &PlantainerCtrlSt{
		*bR,
	}

	return res
}