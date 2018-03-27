package device

import (
	"mevericcore/old/mcdevicerpcmanager_old"
	"mevericcore/old/mcplantainer_old/common"
)

type PlantainerCtrlSt struct {
	mcdevicerpcmanager_old.DeviceRPCCtrlSt
}

func CreateNewPlantainerCtrl() *PlantainerCtrlSt {
	bR := mcdevicerpcmanager_old.CreateNewDeviceRPCCtrl(
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