package mcdashboard

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/gorilla/websocket"
	"mevericcore/mcws"
	"fmt"
	"mevericcore/mccommon"
)

var (
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	WSocketsResources = mcws.CreateNewResourcesManager()
)

func (this *UserController) createAllWSRooms(userId string, userWS *mcws.WSocket) error {
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

func (this *UserController) WSHandler(c echo.Context) error {
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
			if msg.Method != "token" && msg.Method != "authenticate" {
				userWS.SendErrorMsg("Forbidden", msg.Method, 503, msg.Id)
				continue
			}
		}

		if err := WSocketsResources.Handle(msg.Method, userWS, byteMsg); err != nil {
			continue
		}

		//if !appWS.Auth {
		//	//msg := &WsMsgBase{}
		//	//if err := msg.UnmarshalJSON(byteMsg); err != nil {
		//	//	return err
		//	//}
		//	//if msg.Method === "token" {
		//	//	QueueManager.Pub("", msg)
		//	//}
		//	return nil
		//}
		// QueueManager.Pub("ws.msg.receive", msg)
		// if (resourceName) QueueManager.Pub(ws + ".ws.msg.receive", msg)
		//c.Logger().Error(err)
	}

	return c.NoContent(200)
}

