package mcuserrpcmanager

import (
	"mevericcore/mcws"
	"mevericcore/mccommon"
	"github.com/labstack/echo"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserRPCManagerHandleResult struct {
	Resp mccommon.JSONData
	Error mccommon.JSONData
}

type UserRPCManagerHandleResultChannel chan UserRPCManagerHandleResult

type UserRPCManagerSt struct {
	Router *UserRPCRouterSt
	ServerId string
}

func CreateNewUserRPCManagerSt(serverId string) *UserRPCManagerSt {
	return &UserRPCManagerSt{
		CreateNewDeviceRPCRouter(),
		serverId,
	}
}

func (thisR *UserRPCManagerSt) Handle(c UserRPCManagerHandleResultChannel, msg *mccommon.DeviceToServerReqSt) error {
	rpcData := &mccommon.RPCMsg{}
	if err := rpcData.UnmarshalJSON(*msg.Msg); err != nil {
		return thisR.SendRPCErrorRes(c, msg.Protocol, "", msg.DeviceId, 0, err.Error(), 422)
	}

	thisR.Router.Handle(c, rpcData.Method, msg, rpcData)

	return nil
}

func (thisR *UserRPCManagerSt) SendRPCErrorRes(c UserRPCManagerHandleResultChannel, protocol string, methodName string, srcDeviceId string, reqId int, errMessage string, errCode int) error {
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
	c <- UserRPCManagerHandleResult{
		nil,
		data,
	}
	return nil
}

func (thisR *UserRPCManagerSt) SendSuccessResp(c UserRPCManagerHandleResultChannel, msg *mccommon.RPCMsg, result *map[string]interface{}) error {
	data := mccommon.RPCMsg{
		Method: msg.Method,
		Id: msg.Id,
		Src: msg.Dst,
		Dst: msg.Src,
		Result: result,
	}
	c <- UserRPCManagerHandleResult{
		data,
		nil,
	}
	return nil
}

func (thisRPCMan *UserRPCManagerSt) InitRoutes() {
	authG := thisRPCMan.Router.Group("Auth")
	authG.AddHandler("Login", func(req *ReqSt) error {
		tokenMsg := &WsAuthRPCReqSt{}
		if err := tokenMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return err
		}
		if tokenMsg.Args.Login == "" || tokenMsg.Args.Password == "" {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.DeviceId, req.RPCData.Id, "Login and password are required", 503)
		}

		// Create token
		user := new(mccommon.UserModel)

		if err := UsersCollectionManager.FindModelByLogin(tokenMsg.Args.Login, user); err != nil {
			if err == UsersCollectionManager.ErrNotFound {
				return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.DeviceId, req.RPCData.Id, "Invalid email or password", 406)
			} else {
				return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.DeviceId, req.RPCData.Id, "Try again", 406)
			}
		}

		if !user.CheckPasswordHash(tokenMsg.Args.Password) {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.DeviceId, req.RPCData.Id, "Invalid email or password", 406)
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID
		claims["email"] = user.Email
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("secret"))

		if err != nil {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.DeviceId, req.RPCData.Id, "Token creation problem", 503)
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

		req.Channel <- UserRPCManagerHandleResult{
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
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.DeviceId, req.RPCData.Id, "Token is required", 503)
		}
		// Auth user
		t, err := jwt.Parse(authMsg.Args.Token, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.DeviceId, req.RPCData.Id, "Problem with token", 503)
		}

		userTokenId := t.Claims.(jwt.MapClaims)["id"].(string)
		user := new(mccommon.UserModel)

		if err := UsersCollectionManager.FindModelByStringId(userTokenId, user); err != nil {
			return thisRPCMan.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.DeviceId, req.RPCData.Id, "Incorrect token", 503)
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
		userId := req.Msg.DeviceId

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
	//deviceG.AddHandler("List", func(req *mcws.ReqSt) error {
	//	return nil
	//})
	//deviceG.AddHandler("Get", func(req *mcws.ReqSt) error {
	//	return nil
	//})
}