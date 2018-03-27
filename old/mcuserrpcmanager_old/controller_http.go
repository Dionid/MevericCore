package mcuserrpcmanager_old

import (
	"mevericcore/mclibs/mcecho"
	"mevericcore/mclibs/mcws"
	"mevericcore/mclibs/mccommon"
	"github.com/labstack/echo"
	"net/http"
	"fmt"
	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WSHttpControllerSt struct {
	mcecho.ModelControllerBase
}

func (this *WSHttpControllerSt) createAllWSRooms(userId string, userWS *mcws.WSocket) error {
	WSManager.AddWSocketById(userWS)
	WSManager.GetOrAddWSocketRoomWithWSocket(userWS.Id, userWS)

	// Find all devices
	devices := &mccommon.DevicesWithCustomDataListBaseModel{}

	if err := DevicesCollectionManager.FindByOwnerId(userId, devices); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create rooms for them
	for _, dev := range *devices {
		if dev.Shadow.Id != "" {
			WSManager.GetOrAddWSocketRoomWithWSocket(dev.Shadow.Id, userWS)
		}
	}

	return nil
}

func (this *WSHttpControllerSt) WSHandler(c echo.Context) error {
	userId := c.QueryParam("userId")

	ws, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, "WS problem")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "WS problem")
	}

	defer ws.Close()

	userWS := mcws.CreateWSocket(userId, ws)
	this.createAllWSRooms(userId, userWS)

	ws.SetCloseHandler(func(code int, text string) error {
		defer func() {
			if recover() != nil {
				fmt.Println("Recovered")
				return
			}
			fmt.Println("Closed")
		}()
		WSManager.RemoveWSocketById(userId)
		return nil
	})

	for {
		_, byteMsg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			break
		}
		fmt.Println("Receieved: " + string(byteMsg))

		msg := &mcws.WsRPCMsgBaseSt{}
		if err := msg.UnmarshalJSON(byteMsg); err != nil {
			continue
		}

		if !userWS.Authorized {
			if msg.Method != "Auth.Auth" && msg.Method != "Auth.Login" {
				//userWS.SendErrorMsg("Forbidden", msg.Method, 503, msg.Id)
				r := mccommon.RPCMsg{
					Method: msg.Method,
					Id: msg.Id,
					Src: msg.Dst,
					Dst: msg.Src,
					Error: &map[string]interface{}{
						"message": "Forbidden",
						"code": 503,
					},
				}
				if bData, err := r.MarshalJSON(); err != nil {
					userWS.SendMsg([]byte(err.Error()))
				} else {
					userWS.SendMsg(bData)
				}
				continue
			}
		}

		handleMsg := &mccommon.ClientToServerReqSt{
			ClientId:  userId,
			ChannelId: "",
			Protocol:  "WS",
			Msg:       &byteMsg,
		}

		respChan := make(mccommon.ClientToServerHandleResultChannel)

		go func() {
			if err := UserRPCManager.Handle(respChan, handleMsg); err != nil {
				userWS.SendMsg([]byte(err.Error()))
			}
		}()

		for resultSt := range respChan {
			if resultSt.Error != nil {
				if bData, err := resultSt.Error.MarshalJSON(); err != nil {
					userWS.SendMsg([]byte(err.Error()))
				} else {
					userWS.SendMsg(bData)
				}
			}
			if resultSt.Res != nil {
				if bData, err := resultSt.Res.MarshalJSON(); err != nil {
					userWS.SendMsg([]byte(err.Error()))
				} else {
					if msg.Method == "Auth.Auth" {
						userWS.Authorized = true
					}
					userWS.SendMsg(bData)
				}
			}
		}
	}

	return c.NoContent(200)
}
