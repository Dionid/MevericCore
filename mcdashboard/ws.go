package mcdashboard

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/gorilla/websocket"
	"mevericcore/mcws"
	"fmt"
	"mevericcore/mccommon"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

var (
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (this *UserController) createAllWSRooms(userId string, userWS *mcws.WSocket) error {
	WSManager.AddWSocketById(userWS)
	WSManager.GetOrAddWSocketRoomWithWSocket(userWS.Id, userWS)

	// Find all devices
	devices := mccommon.DevicesListBaseModel{}

	if err := DevicesCollectionManager.FindByOwnerId(userId, devices); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create rooms for them
	for _, dev := range devices {
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
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			break
		}
		fmt.Println("Receieved: " + string(msg))

		if !userWS.Authorized {
			msg := &mcws.WsActionMsgBaseSt{}
			if err := msg.UnMarshalJSON(msg); err != nil {
				continue
			}
			if msg.Action == "token" {
				tokenMsg := &WsTokenActionMsgSt{}
				if err := tokenMsg.UnMarshalJSON(msg); err != nil {
					continue
				}
				if tokenMsg.Login == "" || tokenMsg.Password == "" {
					userWS.SendErrorMsg("Login and password are required", msg.Action, 503, msg.RequestId)
					continue
				}
				// Create token
				continue
			}
			if msg.Action == "authenticate" {
				tokenMsg := &WsAuthenticateActionMsgSt{}
				if err := tokenMsg.UnMarshalJSON(msg); err != nil {
					continue
				}
				if tokenMsg.Token == "" {
					userWS.SendErrorMsg("Token is required", msg.Action, 503, msg.RequestId)
					continue
				}
				// Auth user
				t, err := jwt.Parse(tokenMsg.Token, func(t *jwt.Token) (interface{}, error) {
					return []byte("secret"), nil
				})

				if err != nil {
					userWS.SendErrorMsg("Problem with token", msg.Action, 503, msg.RequestId)
				}

				userTokenId := t.Claims.(jwt.MapClaims)["id"].(string)
				user := new(mccommon.UserModel)

				if err := UsersCollectionManager.FindModelByStringId(userTokenId, user); err != nil {
					userWS.SendErrorMsg("Incorrect token", msg.Action, 503, msg.RequestId)
					continue
				}

				 userWS.Authorized = true
				continue
			}

			userWS.SendErrorMsg("Forbidden", msg.Action, 503, msg.RequestId)
		}

		//if !appWS.Auth {
		//	//msg := &WsMsgBase{}
		//	//if err := msg.UnMarshalJSON(msg); err != nil {
		//	//	return err
		//	//}
		//	//if msg.Action === "token" {
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

