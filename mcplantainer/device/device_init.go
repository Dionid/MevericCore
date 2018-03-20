package device

import (
	"mevericcore/mcdevicerpcmanager"
	"mevericcore/mcplantainer/common"
)

type PlantainerCtrlSt struct {
	mcdevicerpcmanager.DeviceRPCCtrlSt
}

func CreateNewPlantainerCtrl(typeName string) *PlantainerCtrlSt {
	bR := mcdevicerpcmanager.CreateNewDeviceRPCCtrl("plantainerServerId", typeName, common.CreateNewPlantainerModelSt, DeviceMQTTManager, SendToUser)
	bR.DevicesCollectionManager = common.PlantainerCollectionManager

	res := &PlantainerCtrlSt{
		*bR,
	}

	return res
}