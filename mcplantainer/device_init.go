package mcplantainer

import "mevericcore/mcdevicerpcmanager"

type PlantainerCtrlSt struct {
	mcdevicerpcmanager.DeviceRPCCtrlSt
}

func CreateNewPlantainerCtrl() *PlantainerCtrlSt {
	bR := mcdevicerpcmanager.CreateNewDeviceRPCCtrl("plantainer")

	res := &PlantainerCtrlSt{
		*bR,
	}

	return res
}