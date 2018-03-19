package mcplantainer

import (
	"mevericcore/mcdevicerpcmanager"
)

type PlantainerCtrlSt struct {
	mcdevicerpcmanager.DeviceRPCCtrlSt
}

func CreateNewPlantainerCtrl(typeName string) *PlantainerCtrlSt {
	bR := mcdevicerpcmanager.CreateNewDeviceRPCCtrl(typeName, CreateNewPlantainerModelSt)
	bR.DevicesCollectionManager = DevicesCollectionManager

	res := &PlantainerCtrlSt{
		*bR,
	}

	return res
}