package device

import (
	"mevericcore/mcdevicerpcmanager"
	"mevericcore/mcplantainer/common"
)

type PlantainerCtrlSt struct {
	mcdevicerpcmanager.DeviceRPCCtrlSt
}

func CreateNewPlantainerCtrl() *PlantainerCtrlSt {
	bR := mcdevicerpcmanager.CreateNewDeviceRPCCtrl(
		PlantainerServerId,
		"plantainer",
		common.PlantainerCollectionManager,
		common.CreateNewPlantainerModelSt,
		InnerRPCMan.SendRPCMsgToUser,
	)

	res := &PlantainerCtrlSt{
		*bR,
	}

	return res
}