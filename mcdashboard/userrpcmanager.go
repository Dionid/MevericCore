package mcdashboard

import (
	"mevericcore/mcuserrpcmanager2"
	"github.com/dgrijalva/jwt-go"
	"time"
	"mevericcore/mcws"
	"mevericcore/mccommunication"
)

var (
	UserRPCManager = mcuserrpcmanager2.New()
)

func InitUserRPCManager() {
	initUserRPCManAuthRoutes()
	initUserRPCManDeviceRoutes()
}

func initUserRPCManDeviceRoutes() {
	//deviceG := UserRPCManager.Router.Group("Device")

}

func initUserRPCManAuthRoutes() {
	authG := UserRPCManager.Router.Group("Auth")
	authG.AddHandler("GetToken", func(req *mccommunication.RPCReqSt) error {
		tokenMsg := &mcws.WsAuthRPCReqSt{}
		if err := tokenMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return err
		}
		if tokenMsg.Args.Login == "" || tokenMsg.Args.Password == "" {
			return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Login and password are required", 503)
		}

		// Create token
		user := new(UserModel)

		if err := usersCollectionManager.FindModelByLogin(tokenMsg.Args.Login, user); err != nil {
			if err == usersCollectionManager.ErrNotFound {
				return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Invalid email or password", 406)
			} else {
				return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Try again", 406)
			}
		}

		if !user.CheckPasswordHash(tokenMsg.Args.Password) {
			return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Invalid email or password", 406)
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID
		claims["email"] = user.Email
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("secret"))

		if err != nil {
			return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Token creation problem", 503)
		}

		// Send success
		rpcResData := &mcws.WsAuthRPCResSt{
			WsRPCMsgBaseSt: mcws.WsRPCMsgBaseSt{
				RPCMsg: mccommunication.RPCMsg{
					Method: tokenMsg.Method,
					Id: tokenMsg.Id,
					Src: tokenMsg.Dst,
					Dst: tokenMsg.Src,
				},
			},
			Result: struct{Token string}{t},
		}

		req.Channel <- mccommunication.ClientToServerHandleResult{
			rpcResData,
			nil,
		}

		return nil
	})

	authG.AddHandler("Auth", func(req *mccommunication.RPCReqSt) error {
		authMsg := &mcws.WsTokenRPCReqSt{}
		if err := authMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return err
		}
		if authMsg.Args.Token == "" {
			return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Token is required", 503)
		}
		// Auth user
		t, err := jwt.Parse(authMsg.Args.Token, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Problem with token", 503)
		}

		userTokenId := t.Claims.(jwt.MapClaims)["id"].(string)
		user := new(UserModel)

		if err := usersCollectionManager.FindModelByStringId(userTokenId, user); err != nil {
			return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Incorrect token", 503)
		}

		// Send success
		res := &map[string]interface{}{"success": true}

		return UserRPCManager.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)
	})
}