package mcplantainer

import (
	"strconv"
	"mevericcore/mcmodules/mclightmodule"
	"github.com/robfig/cron"
	"fmt"
)

func (cr *DeviceCronManagerSt) NewLightModuleSetter() CronSetterFn {
	return func(dId string, c *cron.Cron) error {
		defer func(){
			if recover() != nil {
				fmt.Println("")
				fmt.Println("Recover in NewLightModuleSetter")
				return
			}
		}()
		plantainer := &PlantainerModelSt{}
		if err := plantainerCollectionManager.FindByShadowId(dId, plantainer); err != nil {
			return err
		}

		// . Add intervals cron
		lightModule := plantainer.Shadow.State.Reported.LightModule

		if lightModule.Mode == nil || *lightModule.Mode != mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode] {
			return nil
		}

		//checkTime := strconv.Itoa((*lightModule.LightIntervalsCheckingInterval) / 1000)
		//
		//c.AddFunc("0/" + checkTime + " * * * * *", func(){
		//	plantainer := &PlantainerModelSt{}
		//	if err := plantainerCollectionManager.FindByShadowId(dId, plantainer); err != nil {
		//		return
		//	}
		//
		//	if changed, err := plantainer.CheckAllSystems(); err != nil {
		//		errRPC := NewShadowUpdateRejectedReqRPC(plantainer.Shadow.Id, err.Error(), 500)
		//		innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", errRPC)
		//		innerRPCMan.PublishRPC("User.RPC.Send", errRPC)
		//		return
		//	} else if changed {
		//		if err := plantainerCollectionManager.SaveModel(plantainer); err != nil {
		//			errRPC := NewShadowUpdateRejectedReqRPC(plantainer.Shadow.Id, err.Error(), 500)
		//			innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", errRPC)
		//			innerRPCMan.PublishRPC("User.RPC.Send", errRPC)
		//			return
		//		}
		//		successUpdate := NewShadowUpdateAcceptedReqRPC(
		//			dId,
		//			&plantainer.Shadow,
		//		)
		//
		//		innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", successUpdate)
		//		innerRPCMan.PublishRPC("User.RPC.Send", successUpdate)
		//
		//		plantainer.Shadow.State.FillDelta()
		//
		//		if plantainer.Shadow.State.Delta != nil {
		//			deltaRpc := NewShadowUpdateDeltaReqRPC(dId, &plantainer.Shadow)
		//			innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", deltaRpc)
		//		}
		//	}
		//})

		for _, interval := range *lightModule.LightIntervalsArr {
			fromCrString := "0 " + strconv.Itoa(interval.FromTimeMinutes) + " " + strconv.Itoa(interval.FromTimeHours) + " * * *"
			c.AddFunc(fromCrString, func() {
				plantainer := &PlantainerModelSt{}
				if err := plantainerCollectionManager.FindByShadowId(dId, plantainer); err != nil {
					return
				}
				lightModule := plantainer.Shadow.State.Reported.LightModule

				// Cancel cron action if it happened when mode is in "manual"
				if *lightModule.Mode != mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode] {
					return
				}

				// Change the State (LightTurnedOn)
				plantainer.Shadow.State.Desired.LightModule.LightTurnedOn = &interval.TurnedOn
				plantainer.Shadow.IncrementVersion()

				if err := plantainerCollectionManager.SaveModel(plantainer); err != nil {
					return
				}

				successUpdate := NewShadowUpdateAcceptedReqRPC(
					dId,
					&plantainer.Shadow,
				)

				innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", successUpdate)
				innerRPCMan.PublishRPC("User.RPC.Send", successUpdate)

				plantainer.Shadow.State.FillDelta()

				if plantainer.Shadow.State.Delta != nil {
					deltaRpc := NewShadowUpdateDeltaReqRPC(dId, &plantainer.Shadow)
					innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", deltaRpc)
				}
			})
			toCrStr := "0 " + strconv.Itoa(interval.ToTimeMinutes) + " " + strconv.Itoa(interval.ToTimeHours) + " * * *"
			c.AddFunc(toCrStr, func() {
				plantainer := &PlantainerModelSt{}
				if err := plantainerCollectionManager.FindByShadowId(dId, plantainer); err != nil {
					return
				}
				lightModule := plantainer.Shadow.State.Reported.LightModule
				if *lightModule.Mode != mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode] {
					return
				}

				plantainer.Shadow.State.Desired.LightModule.LightTurnedOn = lightModule.LightIntervalsRestTimeTurnedOn
				plantainer.Shadow.IncrementVersion()

				if err := plantainerCollectionManager.SaveModel(plantainer); err != nil {
					return
				}

				successUpdate := NewShadowUpdateAcceptedReqRPC(
					dId,
					&plantainer.Shadow,
				)

				innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", successUpdate)
				innerRPCMan.PublishRPC("User.RPC.Send", successUpdate)

				plantainer.Shadow.State.FillDelta()

				if plantainer.Shadow.State.Delta != nil {
					deltaRpc := NewShadowUpdateDeltaReqRPC(dId, &plantainer.Shadow)
					innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", deltaRpc)
				}
			})
		}
		return nil
	}
}
