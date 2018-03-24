package mcuserrpcmanager_old

import (
	"mevericcore/mcws"
	"mevericcore/mccommon"
	"github.com/labstack/echo"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
	"errors"
)

type UserRPCManagerSt struct {
	Router *UserRPCRouterSt
	ServerId string
	DeviceCreator mccommon.DeviceCreatorFn
	DevicesListCreator mccommon.DevicesListCreatorFn
	SendToDevice func(msg *mccommon.RPCMsg) error
}

func CreateNewUserRPCManagerSt(serverId string) *UserRPCManagerSt {
	return &UserRPCManagerSt{
		CreateNewDeviceRPCRouter(),
		serverId,
		nil,
		nil,
		func(msg *mccommon.RPCMsg) error {
			if bData, err := msg.MarshalJSON(); err != nil {
				return err
			} else {
				InnerRPCMan.Service.Publish("Device.RPC.Send", bData)
			}

			return nil
		},
	}
}

func (thisR *UserRPCManagerSt) Handle(c mccommon.ClientToServerHandleResultChannel, msg *mccommon.ClientToServerReqSt) error {
	rpcData := &mccommon.RPCMsg{}
	if err := rpcData.UnmarshalJSON(*msg.Msg); err != nil {
		return thisR.SendRPCErrorRes(c, msg.Protocol, "", msg.ClientId, 0, err.Error(), 422)
	}

	thisR.Router.Handle(c, rpcData.Method, msg, rpcData)

	return nil
}

func (thisR *UserRPCManagerSt) SendRPCErrorRes(c mccommon.ClientToServerHandleResultChannel, protocol string, methodName string, srcDeviceId string, reqId int, errMessage string, errCode int) error {
	data := mccommon.RPCMsg{
		Method: methodName,
		Id: reqId,
		Src: thisR.ServerId,
		Dst: srcDeviceId,
		Error: &map[string]interface{}{
			"message": errMessage,
			"code": errCode,
		},
	}
	c <- mccommon.ClientToServerHandleResult{
		nil,
		data,
	}
	return nil
}

func (thisR *UserRPCManagerSt) SendSuccessResp(c mccommon.ClientToServerHandleResultChannel, msg *mccommon.RPCMsg, result *map[string]interface{}) error {
	data := mccommon.RPCMsg{
		Method: msg.Method,
		Id: msg.Id,
		Src: msg.Dst,
		Dst: msg.Src,
		Result: result,
	}
	c <- mccommon.ClientToServerHandleResult{
		data,
		nil,
	}
	return nil
}

func (thisR *UserRPCManagerSt) SendReq(c mccommon.ClientToServerHandleResultChannel, methodName string, srcDeviceId string, reqId int, args *map[string]interface{}) error {
	data := mccommon.RPCMsg{
		Method: methodName,
		Id: reqId,
		Src: thisR.ServerId,
		Dst: srcDeviceId,
		Args: args,
	}
	c <- mccommon.ClientToServerHandleResult{
		data,
		nil,
	}
	return nil
}

func (thisRPCMan *UserRPCManagerSt) Init(deviceCr mccommon.DeviceCreatorFn, devicesLCr mccommon.DevicesListCreatorFn) {
	thisRPCMan.DeviceCreator = deviceCr
	thisRPCMan.DevicesListCreator = devicesLCr
}

func (thisRPCMan *UserRPCManagerSt) InitRoutes() {
	authG := thisRPCMan.Router.Group("Auth")
	authG.AddHandler("Login", func(req *ReqSt) error {
		tokenMsg := &WsAuthRPCReqSt{}
		if err := tokenMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return err
		}
		if tokenMsg.Args.Login == "" || tokenMsg.Args.Password == "" {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "Login and password are required", 503)
		}

		// Create token
		user := new(mccommon.UserModel)

		if err := UsersCollectionManager.FindModelByLogin(tokenMsg.Args.Login, user); err != nil {
			if err == UsersCollectionManager.ErrNotFound {
				return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "Invalid email or password", 406)
			} else {
				return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "Try again", 406)
			}
		}

		if !user.CheckPasswordHash(tokenMsg.Args.Password) {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "Invalid email or password", 406)
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID
		claims["email"] = user.Email
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("secret"))

		if err != nil {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "Token creation problem", 503)
		}

		// Send success
		rpcResData := &WsAuthRPCResSt{
			WsRPCMsgBaseSt: mcws.WsRPCMsgBaseSt{
				RPCMsg: mccommon.RPCMsg{
					Method: tokenMsg.Method,
					Id: tokenMsg.Id,
					Src: tokenMsg.Dst,
					Dst: tokenMsg.Src,
				},
			},
			Result: struct{Token string}{t},
		}

		req.Channel <- mccommon.ClientToServerHandleResult{
			rpcResData,
			nil,
		}

		return nil
	})

	authG.AddHandler("Auth", func(req *ReqSt) error {
		authMsg := &WsTokenRPCReqSt{}
		if err := authMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return err
		}
		if authMsg.Args.Token == "" {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "Token is required", 503)
		}
		// Auth user
		t, err := jwt.Parse(authMsg.Args.Token, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "Problem with token", 503)
		}

		userTokenId := t.Claims.(jwt.MapClaims)["id"].(string)
		user := new(mccommon.UserModel)

		if err := UsersCollectionManager.FindModelByStringId(userTokenId, user); err != nil {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "Incorrect token", 503)
		}

		//req.Ws.Authorized = true

		// Send success
		res := &map[string]interface{}{"success": true}

		return thisRPCMan.SendSuccessResp(req.Channel, req.RPCData, res)
	})
	thisRPCMan.initDeviceResource()
}

func userWsMiddleware(next HandlerFunc) HandlerFunc {
	return func(req *ReqSt) error {
		userId := req.Msg.ClientId

		userM := new(mccommon.UserModel)
		if err := UsersCollectionManager.FindModelByStringId(userId, userM); err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		req.Set("userM", userM)
		if err := next(req); err != nil {
			return err
		}

		return nil
	}
}

func (thisR *UserRPCManagerSt) initDeviceResource() {
	deviceG := thisR.Router.Group("Device")
	deviceG.Use(userWsMiddleware)
	deviceG.AddHandler("List", func(req *ReqSt) error {
		devices := thisR.DevicesListCreator()
		if err := DevicesCollectionManager.FindByOwnerId(req.Msg.ClientId, devices); err != nil {
			return err
		}
		thisR.SendSuccessResp(req.Channel, req.RPCData, &map[string]interface{}{"data": devices})
		return nil
	})
	deviceG.AddHandler("Get", func(req *ReqSt) error {
		device := thisR.DeviceCreator()
		args := req.RPCData.Args.(map[string]interface{})
		deviceId := args["deviceId"].(string)
		if err := DevicesCollectionManager.FindByShadowId(deviceId, device); err != nil {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}
		if isOwner, err := device.IsOwnerStringId(req.Msg.ClientId); !isOwner {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "It's not your device", 403)
		} else if err != nil {
			print(err.Error())
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}
		thisR.SendSuccessResp(req.Channel, req.RPCData, &map[string]interface{}{deviceId: device})
		return nil
	})
	deviceG.AddHandler("Update", func(req *ReqSt) error {
		device := thisR.DeviceCreator()
		args := req.RPCData.Args.(map[string]interface{})
		deviceId := args["deviceId"].(string)
		if err := DevicesCollectionManager.FindByShadowId(deviceId, device); err != nil {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}
		if isOwner, err := device.IsOwnerStringId(req.Msg.ClientId); !isOwner {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "It's not your device", 403)
		} else if err != nil {
			print(err.Error())
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}

		device.Update(&args)

		DevicesCollectionManager.SaveModel(device)

		thisR.SendSuccessResp(req.Channel, req.RPCData, &map[string]interface{}{deviceId: device})
		return nil
	})
	shadowG := deviceG.Group("Shadow")
	shadowG.AddHandler("Get", func(req *ReqSt) error {
		device := thisR.DeviceCreator()
		args := req.RPCData.Args.(map[string]interface{})
		deviceId := args["deviceId"].(string)

		if err := DevicesCollectionManager.FindByShadowId(deviceId, device); err != nil {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}

		if isOwner, err := device.IsOwnerStringId(req.Msg.ClientId); !isOwner {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "It's not your device", 403)
		} else if err != nil {
			print(err.Error())
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}

		state := device.GetShadow().GetState()

		thisR.SendSuccessResp(req.Channel, req.RPCData, &map[string]interface{}{"state": state})

		state.FillDelta()

		if len(state.Delta.State) != 0 {
			data := &mccommon.RPCMsg{
				Method: "Device.Shadow.Delta",
				//Id: reqId,
				Src: thisR.ServerId,
				Dst: deviceId,
				Args: &map[string]interface{}{"state": state.Delta.State, "version": state.Delta.Version, "timestamp": time.Now()},
			}
			thisR.SendToDevice(data)
		}

		return nil
	})
	shadowG.AddHandler("Update", func(req *ReqSt) error {
		device := thisR.DeviceCreator()
		args := req.RPCData.Args.(map[string]interface{})
		deviceId := args["deviceId"].(string)

		if err := DevicesCollectionManager.FindByShadowId(deviceId, device); err != nil {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}

		if isOwner, err := device.IsOwnerStringId(req.Msg.ClientId); !isOwner {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "It's not your device", 403)
		} else if err != nil {
			print(err.Error())
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}

		somethingNew := false
		deviceState := device.GetShadow().GetState()

		updateRpcMsg := &mccommon.RPCWithShadowUpdateMsg{}

		if err := updateRpcMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
		}

		updateData := updateRpcMsg.Args

		device.ActionsOnUpdate(&updateData, DevicesCollectionManager)

		if updateData.State.Reported != nil && updateData.State.Desired != nil {
			deviceState.SetReportedState(updateData.State.Reported)
			deviceState.SetDesiredState(updateData.State.Desired)
			deviceState.IncrementVersion()
			if err := DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
			}
			// PUB /update/accepted with Desire and Reported
			somethingNew = true
		} else if updateData.State.Reported != nil {
			deviceState.SetReportedState(updateData.State.Reported)
			deviceState.IncrementVersion()
			if err := DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
			}
			// PUB /update/accepted with Reported
			somethingNew = true
		} else if updateData.State.Desired != nil {
			if !deviceState.CheckVersion(updateData.Version) {
				// PUB /update/rejected with Desired and Reported
				err := errors.New("version wrong")
				return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
			}
			deviceState.SetDesiredState(updateData.State.Desired)
			deviceState.IncrementVersion()
			if err := DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
			}
			// PUB /update/accepted with Desired
			somethingNew = true
		}

		deviceState.FillDelta()

		if len(deviceState.Delta.State) != 0 {
			data := &mccommon.RPCMsg{
				Method: "Device.Shadow.Delta",
				//Id: reqId,
				Src: thisR.ServerId,
				Dst: deviceId,
				Args: &map[string]interface{}{"state": deviceState.Delta.State, "version": deviceState.Delta.Version, "timestamp": time.Now()},
			}
			thisR.SendToDevice(data)
		}

		if !somethingNew {
			// In this case SetIsActivated haven't been saved
			if err := DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
			}
		}

		thisR.SendSuccessResp(req.Channel, req.RPCData, &map[string]interface{}{"state": deviceState})

		rpcData := &mccommon.RPCMsg{
			Dst: device.GetShadowId(),
			Src: thisR.ServerId,
			Method: "Device.Shadow.Update.Accepted",
			Args: &map[string]interface{}{
				"state": updateData.State,
				"version": deviceState.Metadata.Version,
			},
		}

		thisR.SendToDevice(rpcData)

		return nil
	})
}