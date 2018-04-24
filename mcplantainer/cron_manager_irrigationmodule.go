package mcplantainer

import (
	"github.com/robfig/cron"
	"fmt"
	"mevericcore/mcmodules/mcirrigationmodule"
)

func (cr *DeviceCronManagerSt) NewIrrigationModuleSetter() CronSetterFn {
	return func(dId string, c *cron.Cron) error {
		defer func(){
			if recover() != nil {
				fmt.Println("")
				fmt.Println("Recover in NewIrrigationModuleSetter")
				return
			}
		}()

		plantainer := &PlantainerModelSt{}
		if err := plantainerCollectionManager.FindByShadowId(dId, plantainer); err != nil {
			return err
		}

		irrigationModule := plantainer.Shadow.State.Reported.IrrigationModule

		if irrigationModule.Mode == nil || *irrigationModule.Mode != mcirrigationmodule.IrrigationModuleModeServerIrrigationTimerMode {
			return nil
		}

		//cronTimerString := fmt.Sprintf("*/%x * * * * *", (*irrigationModule.IrrigationTimerEveryXSeconds)/1000)
		//
		//c.AddFunc(cronTimerString, func() {
		//	plantainer := &PlantainerModelSt{}
		//	if err := plantainerCollectionManager.FindByShadowId(dId, plantainer); err != nil {
		//		return
		//	}
		//
		//	irrigationModule := plantainer.Shadow.State.Reported.IrrigationModule
		//
		//	if irrigationModule.Mode == nil || *irrigationModule.Mode != mcirrigationmodule.IrrigationModuleModeServerIrrigationTimerMode {
		//		return
		//	}
		//
		//	t := true
		//	plantainer.Shadow.State.Desired.IrrigationModule.IrrigationTurnedOn = &t
		//	plantainer.Shadow.IncrementVersion()
		//
		//	if err := plantainerCollectionManager.SaveModel(plantainer); err != nil {
		//		return
		//	}
		//
		//	successUpdate := NewShadowUpdateAcceptedReqRPC(
		//		dId,
		//		&plantainer.Shadow,
		//	)
		//
		//	innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", successUpdate)
		//	innerRPCMan.PublishRPC("User.RPC.Send", successUpdate)
		//
		//	plantainer.Shadow.State.FillDelta()
		//
		//	if plantainer.Shadow.State.Delta != nil {
		//		deltaRpc := NewShadowUpdateDeltaReqRPC(dId, &plantainer.Shadow)
		//		innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", deltaRpc)
		//	}
		//})

		return nil
	}
}
