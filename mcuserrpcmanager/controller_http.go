package mcuserrpcmanager

import (
	"mevericcore/mcecho"
	"mevericcore/mcws"
	"mevericcore/mccommon"
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
		WSManager.GetOrAddWSocketRoomWithWSocket(dev.Shadow.Id, userWS)
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
			if msg.Method != "Auth.Token" && msg.Method != "Auth.Authenticate" {
				userWS.SendErrorMsg("Forbidden", msg.Method, 503, msg.Id)
				continue
			}
		}

		handleMsg := &mccommon.DeviceToServerReqSt{
			DeviceId: userId,
			ChannelId: "",
			Protocol: "WS",
			Msg: &byteMsg,
		}

		respChan := make(UserRPCManagerHandleResultChannel)

		if err := UserRPCManager.Handle(respChan, handleMsg); err != nil {
			continue
		}

		for resultSt := range respChan {
			if resultSt.Error != nil {
				if bData, err := resultSt.Error.MarshalJSON(); err != nil {
					userWS.SendMsg([]byte(err.Error()))
				} else {
					userWS.SendMsg(bData)
				}
			}
			if resultSt.Resp != nil {
				if bData, err := resultSt.Resp.MarshalJSON(); err != nil {
					userWS.SendMsg([]byte(err.Error()))
				} else {
					userWS.SendMsg(bData)
				}
			}
		}
	}

	return c.NoContent(200)
}
