package mcdashboard

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/gorilla/websocket"
	"mevericcore/mcws"
)

var (
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (this *UserController) WSHandler(c echo.Context) error {
	userId := c.QueryParam("userId")

	ws, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, "WS problem")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "WS problem")
	}

	defer ws.Close()

	appWS := mcws.CreateWSocket(userId, ws)



	return c.NoContent(200)
}

