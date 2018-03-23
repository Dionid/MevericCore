package mcplantainer

import (
	"mevericcore/mccommunication"
	"mevericcore/mcinnerrpc"
	"fmt"
)

var (
	innerRPCMan = mcinnerrpc.New()
)

func initInnerRPCMan() {
	innerRPCMan.Init()
	innerRPCMan.Service.Subscribe("Devices.Plantainer.RPC.Send", func(msg *mcinnerrpc.Msg) {
		return
	})
	innerRPCMan.Service.Subscribe("User.RPC.Devices.Plantainer.>", func(req *mcinnerrpc.Msg) {
		msg := &mccommunication.ClientToServerRPCReqSt{}
		if err := msg.UnmarshalJSON(req.Data); err != nil {
			fmt.Println("msg.UnmarshalJSON error: " + err.Error())
			return
		}

		respChan := make(mccommunication.ClientToServerHandleResultChannel)

		go func() {
			if err := userRPCManager.Handle(respChan, msg); err != nil {
				data := &mccommunication.RPCMsg{
					Method: msg.RPCMsg.Method,
					Id: msg.RPCMsg.Id,
					Src: msg.RPCMsg.Dst,
					Dst: msg.RPCMsg.Src,
					Error: &map[string]interface{}{
						"message": err.Error(),
						"code": 500,
					},
				}
				if bData, err := data.MarshalJSON(); err != nil {
					data := &mccommunication.RPCMsg{
						Method: msg.RPCMsg.Method,
						Id: msg.RPCMsg.Id,
						Src: msg.RPCMsg.Dst,
						Dst: msg.RPCMsg.Src,
						Error: &map[string]interface{}{
							"message": "Marshaling error problem",
							"code": 500,
						},
					}
					ebData, _ := data.MarshalJSON()
					innerRPCMan.Service.Publish("User.RPC.Send", ebData)
					print(bData)
					return
				} else {
					return
					//userWS.SendMsg(bData)
				}
			}
		}()

		for resultSt := range respChan {
			if resultSt.Error != nil {
				if bData, err := resultSt.Error.MarshalJSON(); err != nil {
					data := &mccommunication.RPCMsg{
						Method: msg.RPCMsg.Method,
						Id: msg.RPCMsg.Id,
						Src: msg.RPCMsg.Dst,
						Dst: msg.RPCMsg.Src,
						Error: &map[string]interface{}{
							"message": "Marshaling error problem",
							"code": 500,
						},
					}
					ebData, _ := data.MarshalJSON()
					innerRPCMan.Service.Publish("User.RPC.Send", ebData)
				} else {
					innerRPCMan.Service.Publish("User.RPC.Send", bData)
				}
			}
			if resultSt.Res != nil {
				if bData, err := resultSt.Res.MarshalJSON(); err != nil {
					data := &mccommunication.RPCMsg{
						Method: msg.RPCMsg.Method,
						Id: msg.RPCMsg.Id,
						Src: msg.RPCMsg.Dst,
						Dst: msg.RPCMsg.Src,
						Error: &map[string]interface{}{
							"message": "Marshaling error problem",
							"code": 500,
						},
					}
					ebData, _ := data.MarshalJSON()
					innerRPCMan.Service.Publish("User.RPC.Send", ebData)
				} else {
					innerRPCMan.Service.Publish("User.RPC.Send", bData)
				}
			}
		}
	})
}
