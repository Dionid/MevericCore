package mcplantainer

import (
	"mevericcore/mclibs/mccommunication"
	"mevericcore/mclibs/mcinnerrpc"
	"fmt"
)

var (
	innerRPCMan = mcinnerrpc.New()
)

func initInnerRPCMan() {
	innerRPCMan.Init()
	innerRPCMan.Service.Subscribe("Plantainer.Cron.RPC", func(req *mcinnerrpc.Msg){
		msg := &mccommunication.RPCMsg{}
		if err := msg.UnmarshalJSON(req.Data); err != nil {
			fmt.Println("msg.UnmarshalJSON error: " + err.Error())
			return
		}

		args := msg.Args.(map[string]interface{})
		deviceId := args["deviceId"].(string)

		respChan := make(mccommunication.ClientToServerHandleResultChannel)

		go func() {
			if err := cronRPCMan.Handle(respChan, msg, &req.Data); err != nil {
				return
			}
		}()

		for resultSt := range respChan {
			if resultSt.Error != nil {
				if bData, err := resultSt.Error.MarshalJSON(); err != nil {
					rpcData := &mccommunication.RPCMsg{
						Method: msg.Method,
						Id: msg.Id, // ToDo: Add ReqId
						Src: PlantainerServerId,
						Dst: deviceId,
						Error: &map[string]interface{}{
							"message": "Marshaling error problem",
							"code": 500,
						},
					}
					innerRPCMan.PublishRPC("User.RPC.Send", rpcData)
					innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", rpcData)
					return
				} else {
					innerRPCMan.Service.Publish("User.RPC.Send", bData)
					innerRPCMan.Service.Publish("Plantainer.Device.RPC.Send", bData)
				}
			}
			if resultSt.Res != nil {
				if bData, err := resultSt.Res.MarshalJSON(); err != nil {
					rpcData := &mccommunication.RPCMsg{
						Method: msg.Method,
						Id: msg.Id, // ToDo: Add ReqId
						Src: PlantainerServerId,
						Dst: deviceId,
						Error: &map[string]interface{}{
							"message": "Marshaling error problem",
							"code": 500,
						},
					}
					innerRPCMan.PublishRPC("User.RPC.Send", rpcData)
					innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", rpcData)
					return
				} else {
					innerRPCMan.Service.Publish("User.RPC.Send", bData)
					innerRPCMan.Service.Publish("Plantainer.Device.RPC.Send", bData)
				}
			}
		}
	})

	innerRPCMan.Service.Subscribe("Plantainer.Device.RPC.Send", func(msg *mcinnerrpc.Msg) {
		rpcMsg := &mccommunication.RPCMsg{}
		if err := rpcMsg.UnmarshalJSON(msg.Data); err != nil {
			return
		}
		deviceMQTTMan.Publish(rpcMsg.Dst+"/rpc", msg.Data)
	})

	// About Nats at all:
	// There are . variants of subjects:
	// 1. Call some service. Like "Plantainer.Corn.RPC".
	// 2. Generate some event. Like "User.RPC.Received"
	// 3. ...

	// How this convention works:
	// Some app, that receives requests from different sources (Http, WS, MQTT, etc.)
	// generates event "{Source}.RPC.Received", that is common for EVERY message.
	// This subject is used for logging messages and thats all.
	// For more specific we added RPC.Method itself to subject, for you to catch it anywhere is your architecture.
	// Example: "{Source}.RPC.Received.{RPCMsg.Method}"

	//innerRPCMan.Service.Subscribe("Device.RPC.Received.Plantainer>", func(req *mcinnerrpc.Msg) {})
	//innerRPCMan.Service.Subscribe("Cloud.RPC.Received.Plantainer>", func(req *mcinnerrpc.Msg) {})
	innerRPCMan.Service.Subscribe("User.RPC.Received.Plantainer>", func(req *mcinnerrpc.Msg) {
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
